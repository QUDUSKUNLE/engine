package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocketManager manages all WebSocket connections and notifications
type (
	WebSocketManager struct {
		clients    map[string]*Client // userID -> Client
		register   chan *Client
		unregister chan *Client
		broadcast  chan *NotificationMessage
		mutex      sync.RWMutex
	}
	// Client represents a WebSocket client connection
	Client struct {
		ID       string
		UserID   string
		Conn     *websocket.Conn
		Send     chan *NotificationMessage
		Manager  *WebSocketManager
		LastSeen time.Time
	}
	// NotificationMessage represents a real-time notification
	NotificationMessage struct {
		Type      string                 `json:"type"`
		UserID    string                 `json:"user_id"`
		Title     string                 `json:"title"`
		Message   string                 `json:"message"`
		Data      map[string]interface{} `json:"data,omitempty"`
		Timestamp time.Time              `json:"timestamp"`
		ID        string                 `json:"id"`
	}
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *NotificationMessage),
	}
}

// Start begins the WebSocket manager's main loop
func (manager *WebSocketManager) Start() {
	go manager.run()
	utils.Info("WebSocket manager started")
}

// run handles the main WebSocket manager loop
func (manager *WebSocketManager) run() {
	ticker := time.NewTicker(30 * time.Second) // Ping interval
	defer ticker.Stop()

	for {
		select {
		case client := <-manager.register:
			manager.registerClient(client)

		case client := <-manager.unregister:
			manager.unregisterClient(client)

		case message := <-manager.broadcast:
			manager.broadcastMessage(message)

		case <-ticker.C:
			manager.pingClients()
		}
	}
}

// registerClient adds a new client to the manager
func (manager *WebSocketManager) registerClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	// Close existing connection if user is already connected
	if existingClient, exists := manager.clients[client.UserID]; exists {
		close(existingClient.Send)
		existingClient.Conn.Close()
	}

	manager.clients[client.UserID] = client
	utils.Info("WebSocket client registered", utils.LogField{Key: "user_id", Value: client.UserID}, utils.LogField{Key: "client_id", Value: client.ID})

	// Send welcome message
	welcomeMsg := &NotificationMessage{
		Type:      "connection_established",
		UserID:    client.UserID,
		Title:     "Connected",
		Message:   "Real-time notifications enabled",
		Timestamp: time.Now(),
		ID:        uuid.New().String(),
	}

	select {
	case client.Send <- welcomeMsg:
	default:
		close(client.Send)
		delete(manager.clients, client.UserID)
	}
}

// unregisterClient removes a client from the manager
func (manager *WebSocketManager) unregisterClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if _, exists := manager.clients[client.UserID]; exists {
		delete(manager.clients, client.UserID)
		close(client.Send)
		utils.Info("WebSocket client unregistered", utils.LogField{Key: "user_id", Value: client.UserID})
	}
}

// broadcastMessage sends a message to the appropriate client
func (manager *WebSocketManager) broadcastMessage(message *NotificationMessage) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if client, exists := manager.clients[message.UserID]; exists {
		select {
		case client.Send <- message:
		default:
			// Client's send channel is full, remove the client
			close(client.Send)
			delete(manager.clients, message.UserID)
		}
	}
}

// pingClients sends ping messages to all connected clients
func (manager *WebSocketManager) pingClients() {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for userID, client := range manager.clients {
		if time.Since(client.LastSeen) > 60*time.Second {
			// Client hasn't been seen for too long, remove it
			close(client.Send)
			delete(manager.clients, userID)
			continue
		}

		if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			close(client.Send)
			delete(manager.clients, userID)
		}
	}
}

// SendNotification sends a notification to a specific user
func (manager *WebSocketManager) SendNotification(userID string, notification *domain.Notification) {
	message := &NotificationMessage{
		Type:      string(notification.Type),
		UserID:    userID,
		Title:     notification.Title,
		Message:   notification.Message,
		Data:      notification.Metadata,
		Timestamp: notification.CreatedAt,
		ID:        notification.ID.String(),
	}

	select {
	case manager.broadcast <- message:
	default:
		utils.Warn("Failed to queue notification for broadcast", utils.LogField{Key: "user_id", Value: userID})
	}
}

// SendCustomNotification sends a custom notification message
func (manager *WebSocketManager) SendCustomNotification(userID, notificationType, title, message string, data map[string]interface{}) {
	notification := &NotificationMessage{
		Type:      notificationType,
		UserID:    userID,
		Title:     title,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		ID:        uuid.New().String(),
	}

	select {
	case manager.broadcast <- notification:
	default:
		utils.Warn("Failed to queue custom notification for broadcast", utils.LogField{Key: "user_id", Value: userID})
	}
}

// GetConnectedClients returns the number of connected clients
func (manager *WebSocketManager) GetConnectedClients() int {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return len(manager.clients)
}

// IsUserConnected checks if a user is currently connected
func (manager *WebSocketManager) IsUserConnected(userID string) bool {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	_, exists := manager.clients[userID]
	return exists
}

// HandleWebSocket handles WebSocket connection upgrades
func (manager *WebSocketManager) HandleWebSocket(c echo.Context) error {
	// Extract user ID from JWT token or query parameter
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user_id is required",
		})
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		utils.Error("Failed to upgrade WebSocket connection", utils.LogField{Key: "error", Value: err.Error()})
		return err
	}

	// Create new client
	client := &Client{
		ID:       uuid.New().String(),
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan *NotificationMessage, 256),
		Manager:  manager,
		LastSeen: time.Now(),
	}

	// Register client
	manager.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()

	return nil
}

// readPump handles reading messages from the WebSocket connection
func (client *Client) readPump() {
	defer func() {
		client.Manager.unregister <- client
		client.Conn.Close()
	}()

	// Set read deadline and pong handler
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.LastSeen = time.Now()
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// Read message from client
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Error("WebSocket read error", utils.LogField{Key: "error", Value: err.Error()})
			}
			break
		}

		client.LastSeen = time.Now()

		// Handle client messages (e.g., mark notifications as read)
		client.handleMessage(message)
	}
}

// writePump handles writing messages to the WebSocket connection
func (client *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(message); err != nil {
				utils.Error("WebSocket write error", utils.LogField{Key: "error", Value: err.Error()})
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes messages received from the client
func (client *Client) handleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		utils.Error("Failed to parse client message", utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "ping":
		// Respond with pong
		response := &NotificationMessage{
			Type:      "pong",
			UserID:    client.UserID,
			Timestamp: time.Now(),
			ID:        uuid.New().String(),
		}
		select {
		case client.Send <- response:
		default:
		}

	case "mark_read":
		// Handle marking notifications as read
		if notificationID, ok := msg["notification_id"].(string); ok {
			utils.Info("Notification marked as read via WebSocket",
				utils.LogField{Key: "user_id", Value: client.UserID},
				utils.LogField{Key: "notification_id", Value: notificationID})
		}
	}
}
