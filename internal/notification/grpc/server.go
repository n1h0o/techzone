package grpc

import (
	"context"
	"techzone/internal/notification/pb"
	"techzone/internal/notification/service"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer

	notificationService *service.NotificationService
}

func NewServer(
	notificationService *service.NotificationService,
) *Server {
	return &Server{
		notificationService: notificationService,
	}
}

func Register(
	grpcServer *grpc.Server,
	server *Server,
) {
	pb.RegisterNotificationServiceServer(
		grpcServer,
		server,
	)
}

func (s *Server) GetNotifications(
	ctx context.Context,
	req *pb.GetNotificationsRequest,
) (*pb.GetNotificationsResponse, error) {
	notifications, err := s.notificationService.GetNotifications(
		ctx,
		req.GetUserId(),
	)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetNotificationsResponse{
		Notifications: make([]*pb.Notification, 0, len(notifications)),
	}

	for _, n := range notifications {
		resp.Notifications = append(resp.Notifications, &pb.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			OrderId:   n.OrderID,
			Message:   n.Message,
			CreatedAt: timestamppb.New(n.CreatedAt),
		})
	}
	return resp, nil
}
