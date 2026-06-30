package service

import "context"

type NotificationCreator interface {
	CreateNotification(
		ctx context.Context,
		userID int64,
		orderID int64,
		message string,
	) error
}
