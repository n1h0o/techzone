package model

import "time"

type Notification struct {
	ID        int64
	UserID    int64
	OrderID   int64
	Message   string
	CreatedAt time.Time
}
