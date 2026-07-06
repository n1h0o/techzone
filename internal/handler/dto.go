package handler

import "techzone/internal/model"

type MessageResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type MeResponse struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

type ProductResponse struct {
	ID int64 `json:"id"`
}

type ProductStatusInput struct {
	IsActive bool `json:"is_active"`
}

type CartResponse struct {
	Items []model.CartItemInfo `json:"items"`
}

type OrderCreateResponse struct {
	OrderID int64 `json:"order_id"`
}

type OrdersResponse struct {
	Orders []model.OrderInfo `json:"orders"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" example:"processing"`
}

type NotificationsResponse struct {
	Notifications []model.Notification `json:"notifications"`
}

type PaymentRequest struct {
	OrderID int64 `json:"order_id" example:"15"`
}
