package repository

import (
	"context"
	"techzone/internal/model"
)

type NotificationRepository struct {
	db DBTX
}

func NewNotificationRepository(
	db DBTX,
) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Create(
	ctx context.Context,
	notification *model.Notification,
) error {
	query :=
		`
	INSERT INTO notifications(
	user_id,
	order_id,
	message
	)
	VALUES($1,$2,$3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		notification.UserID,
		notification.OrderID,
		notification.Message,
	)
	return err
}

func (r *NotificationRepository) GetNotifications(
	ctx context.Context,
	userID int64,
) ([]model.Notification, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
		id,
		user_id,
		order_id,
		message,
		created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notifications := make([]model.Notification, 0)

	for rows.Next() {
		var notification model.Notification

		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.OrderID,
			&notification.Message,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}
