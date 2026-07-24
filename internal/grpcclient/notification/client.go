package notification

import (
	pb "techzone/internal/notification/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	pb.NotificationServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:                      conn,
		NotificationServiceClient: pb.NewNotificationServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
