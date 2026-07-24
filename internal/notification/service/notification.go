package service

import (
	"context"
	"techzone/internal/notification/model"
)

type NotificationRepository interface {
	Create(
		ctx context.Context,
		notification *model.Notification,
	) error

	GetNotifications(
		ctx context.Context,
		userId int64,
	) ([]model.Notification, error)
}

type NotificationService struct {
	notificationRepo NotificationRepository
}

func NewNotificationService(
	notificationRepo NotificationRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	userID int64,
	orderID int64,
	message string,
) error {
	return s.notificationRepo.Create(
		ctx,
		&model.Notification{
			UserID:  userID,
			OrderID: orderID,
			Message: message,
		},
	)
}

func (s *NotificationService) GetNotifications(
	ctx context.Context,
	userID int64,
) ([]model.Notification, error) {
	return s.notificationRepo.GetNotifications(
		ctx,
		userID,
	)
}
