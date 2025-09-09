package api_check_suite

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)
func TestWebSocketConnection(t *testing.T) {
	// WebSocket server URL
	// WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "host:7556", Path: "/v1/ws/notifications"}
	q := u.Query()
	q.Set("user_id", "test-user-123")
	u.RawQuery = q.Encode()

	fmt.Printf("Connecting to %s\n", u.String())

	// Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer c.Close()

	// Channel to handle interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel to receive messages
	done := make(chan struct{})

	// Goroutine to read messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Printf("📬 Received: %s\n", message)
		}
	}()

	// Send a ping message every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Println("✅ Connected! Waiting for notifications...")
	fmt.Println("💡 You can now send test notifications from another terminal:")
	fmt.Println("   curl -X POST http://localhost:7556/v1/ws/test-notification \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"user_id\":\"test-user-123\",\"title\":\"Test\",\"message\":\"Hello WebSocket!\"}'")
	fmt.Println("\n🛑 Press Ctrl+C to exit")

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Send ping
			err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")

			// Send close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
