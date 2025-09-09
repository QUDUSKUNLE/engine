package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
)

func (r *Repository) CreateNotification(ctx context.Context, arg db.CreateNotificationParams) (*db.Notification, error) {
	return r.database.CreateNotification(ctx, arg)
}

func (r *Repository) GetNotifications(ctx context.Context, arg db.GetUserNotificationsParams) ([]*db.Notification, error) {
	return r.database.GetUserNotifications(ctx, arg)
}

func (r *Repository) MarkNotificationRead(ctx context.Context, arg db.MarkAsReadParams) (*db.Notification, error) {
	return r.database.MarkAsRead(ctx, arg)
}

func (r *Repository) MarkNotificationReadAll(ctx context.Context, arg string) error {
	return r.database.MarkAllAsRead(ctx, arg)
}

var _ ports.NotificationRepository = (*Repository)(nil)
