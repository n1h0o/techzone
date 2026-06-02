package model

import "time"

type User struct {
	ID           int64
	Login        string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}
