package services

import (
	"github.com/medivue/core/utils"
)

func (service *ServicesHandler) emailGoroutine(data interface{}, email, subject, template string) {
	err := service.notificationService.SendEmail(email, subject, template, data)
	if err != nil {
		utils.Error("Failed to send verification email",
			utils.LogField{Key: "error", Value: err.Error()})
	}
	utils.Info("Email notification sent successfully",
		utils.LogField{Key: "user_id"},
		utils.LogField{Key: "user_type"})
}
