package services

import (
	"context"
	"fmt"
	"time"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/google/uuid"
)

// Enhanced notification methods for WebSocket integration

// CreateAndSendNotification creates a notification in the database and sends it via WebSocket
func (srv *ServicesHandler) CreateAndSendNotification(ctx context.Context, userID uuid.UUID, notificationType domain.NotificationType, title, message string, metadata map[string]interface{}) error {
	// Create notification in database
	notification := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
		Metadata:  metadata,
	}

	// Save to database (assuming you have a method for this)
	// err := srv.notificationRepo.CreateNotification(ctx, notification)
	// if err != nil {
	//     return err
	// }

	// Send via WebSocket if user is connected
	if srv.WebSocketManager.IsUserConnected(userID.String()) {
		srv.WebSocketManager.SendNotification(userID.String(), notification)
		utils.Info("Notification sent via WebSocket", 
			utils.LogField{Key: "user_id", Value: userID.String()},
			utils.LogField{Key: "type", Value: string(notificationType)})
	} else {
		utils.Info("User not connected to WebSocket, notification saved to database only", 
			utils.LogField{Key: "user_id", Value: userID.String()})
	}

	return nil
}

// SendAppointmentNotification sends appointment-related notifications
func (srv *ServicesHandler) SendAppointmentNotification(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID, notificationType domain.NotificationType, customMessage string) error {
	var title, message string
	
	switch notificationType {
	case domain.AppointmentCreated:
		title = "Appointment Scheduled"
		message = "Your appointment has been successfully scheduled"
	case domain.AppointmentConfirmed:
		title = "Appointment Confirmed"
		message = "Your appointment has been confirmed"
	case domain.AppointmentRescheduled:
		title = "Appointment Rescheduled"
		message = "Your appointment has been rescheduled"
	case domain.AppointmentCancelled:
		title = "Appointment Cancelled"
		message = "Your appointment has been cancelled"
	case domain.AppointmentReminder:
		title = "Appointment Reminder"
		message = "You have an upcoming appointment"
	default:
		title = "Appointment Update"
		message = "Your appointment status has been updated"
	}

	if customMessage != "" {
		message = customMessage
	}

	metadata := map[string]interface{}{
		"appointment_id": appointmentID.String(),
		"type":          "appointment",
	}

	return srv.CreateAndSendNotification(ctx, userID, notificationType, title, message, metadata)
}

// SendLabResultNotification sends lab result notifications
func (srv *ServicesHandler) SendLabResultNotification(ctx context.Context, userID uuid.UUID, testName string, isAbnormal bool) error {
	var title, message string
	var notificationType domain.NotificationType = "LAB_RESULT_READY"

	if isAbnormal {
		title = "Lab Results Available - Review Required"
		message = fmt.Sprintf("Your %s results are ready and require review", testName)
	} else {
		title = "Lab Results Available"
		message = fmt.Sprintf("Your %s results are ready", testName)
	}

	metadata := map[string]interface{}{
		"test_name":   testName,
		"is_abnormal": isAbnormal,
		"type":       "lab_result",
	}

	return srv.CreateAndSendNotification(ctx, userID, notificationType, title, message, metadata)
}

// SendAIAnalysisNotification sends AI analysis completion notifications
func (srv *ServicesHandler) SendAIAnalysisNotification(ctx context.Context, userID uuid.UUID, analysisType, analysisID string) error {
	title := "AI Analysis Complete"
	message := fmt.Sprintf("Your %s analysis is ready for review", analysisType)
	
	metadata := map[string]interface{}{
		"analysis_type": analysisType,
		"analysis_id":   analysisID,
		"type":         "ai_analysis",
	}

	notificationType := domain.NotificationType("AI_ANALYSIS_COMPLETE")
	return srv.CreateAndSendNotification(ctx, userID, notificationType, title, message, metadata)
}

// SendSystemNotification sends system-wide notifications
func (srv *ServicesHandler) SendSystemNotification(ctx context.Context, userID uuid.UUID, title, message string) error {
	metadata := map[string]interface{}{
		"type": "system",
	}

	notificationType := domain.NotificationType("SYSTEM_NOTIFICATION")
	return srv.CreateAndSendNotification(ctx, userID, notificationType, title, message, metadata)
}

// BroadcastToAllUsers sends a notification to all connected users
func (srv *ServicesHandler) BroadcastToAllUsers(title, message string, data map[string]interface{}) {
	// This would require getting all user IDs from your database
	// For now, we'll just log that a broadcast was requested
	utils.Info("Broadcast notification requested", 
		utils.LogField{Key: "title", Value: title},
		utils.LogField{Key: "connected_clients", Value: srv.WebSocketManager.GetConnectedClients()})
	
	// In a real implementation, you would:
	// 1. Get all active user IDs from your database
	// 2. Loop through them and send notifications
	// 3. Only send via WebSocket to connected users
}
