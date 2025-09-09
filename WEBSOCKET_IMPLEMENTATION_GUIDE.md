# WebSocket Real-Time Notifications Implementation Guide

## 🚀 Overview

Your Diagnoxix now has a comprehensive WebSocket-based real-time notification system that enables instant delivery of notifications to connected users. This system is perfect for healthcare applications where timely communication is critical.

## ✅ What's Implemented

### 1. WebSocket Manager (`core/services/websocket/index.go`)
- **Connection Management**: Handles multiple concurrent WebSocket connections
- **User Mapping**: Maps user IDs to WebSocket connections
- **Message Broadcasting**: Sends notifications to specific users
- **Connection Health**: Automatic ping/pong and connection cleanup
- **Graceful Handling**: Proper connection registration/unregistration

### 2. WebSocket Handlers (`adapters/handlers/websocket.handler.go`)
- **Connection Endpoint**: `/v1/ws/notifications` for establishing WebSocket connections
- **Statistics Endpoint**: `/v1/ws/stats` for monitoring connection status
- **Test Notifications**: `/v1/ws/test-notification` for testing the system

### 3. Enhanced Notification Service (`core/services/notification.websocket.service.go`)
- **Integrated Delivery**: Combines database storage with real-time WebSocket delivery
- **Appointment Notifications**: Specialized methods for appointment-related notifications
- **Lab Result Notifications**: Notifications for lab results and AI analysis
- **System Notifications**: Broadcast capabilities for system-wide messages

### 4. WebSocket Routes (`adapters/routes/websocket.go`)
- **RESTful Structure**: Clean URL structure under `/v1/ws/`
- **Proper HTTP Methods**: GET for connections, POST for test notifications
- **Route Registration**: Integrated with your existing routing system

## 🔌 WebSocket Endpoints

### 1. Establish WebSocket Connection
```
GET /v1/ws/notifications?user_id={user_id}
```
- **Purpose**: Upgrade HTTP connection to WebSocket
- **Parameters**: `user_id` (required) - User identifier
- **Response**: WebSocket connection established (HTTP 101)

### 2. Get Connection Statistics
```
GET /v1/ws/stats
```
- **Purpose**: Monitor WebSocket connections
- **Response**: Number of connected clients and system status

### 3. Send Test Notification
```
POST /v1/ws/test-notification
```
- **Purpose**: Send test notifications for development/testing
- **Body**: JSON with user_id, title, message, and optional data

## 📱 Client Integration

### JavaScript WebSocket Client
```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:7556/v1/ws/notifications?user_id=user-123');

// Handle connection events
ws.onopen = function(event) {
    console.log('Connected to WebSocket');
};

ws.onmessage = function(event) {
    const notification = JSON.parse(event.data);
    console.log('Received notification:', notification);
    // Display notification in UI
    displayNotification(notification);
};

ws.onclose = function(event) {
    console.log('WebSocket connection closed');
};

// Send ping to keep connection alive
setInterval(() => {
    if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'ping' }));
    }
}, 30000);
```

### React Hook Example
```javascript
import { useState, useEffect, useRef } from 'react';

export const useWebSocketNotifications = (userId) => {
    const [notifications, setNotifications] = useState([]);
    const [connected, setConnected] = useState(false);
    const ws = useRef(null);

    useEffect(() => {
        if (!userId) return;

        const wsUrl = `ws://localhost:7556/v1/ws/notifications?user_id=${userId}`;
        ws.current = new WebSocket(wsUrl);

        ws.current.onopen = () => {
            setConnected(true);
        };

        ws.current.onmessage = (event) => {
            const notification = JSON.parse(event.data);
            setNotifications(prev => [notification, ...prev.slice(0, 9)]);
        };

        ws.current.onclose = () => {
            setConnected(false);
        };

        return () => {
            ws.current?.close();
        };
    }, [userId]);

    return { notifications, connected };
};
```

## 🔧 Backend Integration

### Sending Notifications from Your Services

```go
// Send appointment notification
err := srv.SendAppointmentNotification(
    ctx,
    userID,
    appointmentID,
    domain.AppointmentConfirmed,
    "Your appointment has been confirmed for tomorrow at 2:00 PM",
)

// Send lab result notification
err := srv.SendLabResultNotification(
    ctx,
    userID,
    "Complete Blood Count",
    false, // not abnormal
)

// Send AI analysis notification
err := srv.SendAIAnalysisNotification(
    ctx,
    userID,
    "symptom_analysis",
    analysisID,
)

// Send custom notification
err := srv.CreateAndSendNotification(
    ctx,
    userID,
    domain.NotificationType("CUSTOM_TYPE"),
    "Custom Title",
    "Custom message",
    map[string]interface{}{
        "custom_data": "value",
    },
)
```

### Integration with Existing Services

```go
// In your appointment service
func (srv *ServicesHandler) ConfirmAppointment(ctx context.Context, appointmentID uuid.UUID) error {
    // ... existing appointment confirmation logic ...
    
    // Send real-time notification
    err := srv.SendAppointmentNotification(
        ctx,
        appointment.PatientID,
        appointmentID,
        domain.AppointmentConfirmed,
        "",
    )
    if err != nil {
        utils.Warn("Failed to send appointment notification", 
            utils.LogField{Key: "error", Value: err.Error()})
    }
    
    return nil
}
```

## 📋 Notification Message Format

```json
{
    "type": "appointment_confirmed",
    "user_id": "user-123",
    "title": "Appointment Confirmed",
    "message": "Your appointment has been confirmed for tomorrow at 2:00 PM",
    "data": {
        "appointment_id": "appointment-456",
        "type": "appointment"
    },
    "timestamp": "2024-01-15T10:30:00Z",
    "id": "notification-789"
}
```

## 🔒 Security Features

### Authentication
- User ID validation on connection
- JWT token support (can be extended)
- Origin checking (configurable)

### Connection Management
- Automatic cleanup of stale connections
- Rate limiting through existing middleware
- Proper error handling and logging

### Data Validation
- Input validation on all endpoints
- Structured error responses
- Secure WebSocket upgrade process

## 📊 Monitoring and Debugging

### Connection Statistics
```bash
curl http://localhost:7556/v1/ws/stats
```

### Logs
- Connection events are logged with user IDs
- Message delivery status is tracked
- Error conditions are properly logged

### Health Checks
- Automatic ping/pong for connection health
- Cleanup of inactive connections
- Connection timeout handling

## 🧪 Testing

### 1. Test WebSocket Connection
```bash
# Start your server
go run main.go

# Open the test client
open examples/websocket_client.html
```

### 2. Send Test Notification
```bash
curl -X POST http://localhost:7556/v1/ws/test-notification \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "titTest Notification",
    "message": "This is a test message"
  }'
```

### 3. Check Connection Stats
```bash
curl http://localhost:7556/v1/ws/stats
```
🚀 Production Considerations

### Scaling
- Consider using Redis for multi-instance deployments
- Implement connection pooling for high traffic
- Add load balancer WebSocket support

### Performance
- Monitor connection counts and memory usage
- Implement message queuing for offline users
- Add connection rate limiting

### Reliability
- Implement message acknowledgments
- Add retry logic for failed deliveries
- Consider message persistence for critical notifications

## 🔄 Integration with Existing Features

### With AI Services
```go
// After AI analysis completion
interpretation, err := ai.InterpretLabResults(ctx, labTest)
if err == nil {
    srv.SendAIAnalysisNotification(ctx, userID, "lab_interpretation", interpretation.ID)
}
```

### With Appointment System
```go
// After appointment creation
appointment, err := srv.CreateAppointment(ctx, appointmentData)
if err == nil {
    srv.SendAppointmentNotification(ctx, appointment.PatientID, appointment.ID, domain.AppointmentCreated, "")
}
```

### With Notification System
- Seamlessly integrates with existing notification database
- Maintains backward compatibility
- Enhances existing email/SMS notifications

## 📈 Future Enhancements

1. **Message Acknowledgments**: Ensure message delivery
2. **Offline Message Queue**: Store messages for offline users
3. **Push Notifications**: Integration with mobile push services
4. **Message History**: WebSocket-based message history retrieval
5. **User Presence**: Show online/offline status
6. **Typing Indicators**: For chat-like features
7. **File Sharing**: Send files through WebSocket
8. **Video Call Integration**: WebRTC signaling support

Your Diagnoxix application now has enterprise-grade real-time notification capabilities that will significantly enhance user experience and engagement! 🎉
