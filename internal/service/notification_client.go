package service

import (
	"context"
	notificationgrpc "techzone/internal/grpcclient/notification"
	"techzone/internal/notification/pb"
)

type NotificationClientService struct {
	client *notificationgrpc.Client
}

func NewNotificationClientService(
	client *notificationgrpc.Client,
) *NotificationClientService {
	return &NotificationClientService{
		client: client,
	}
}

func (s *NotificationClientService) GetNotifications(
	ctx context.Context,
	userID int64,
) ([]*pb.Notification, error) {
	resp, err := s.client.GetNotifications(
		ctx,
		&pb.GetNotificationsRequest{
			UserId: userID,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Notifications, nil
}
