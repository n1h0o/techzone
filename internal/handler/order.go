package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"techzone/internal/middleware"
	"techzone/internal/model"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type OrderHandler struct {
	orderService     *service.OrderService
	notificationPool *service.NotificationWorkerPool
}

func NewOrderHandler(
	orderService *service.OrderService,
	pool *service.NotificationWorkerPool,
) *OrderHandler {
	return &OrderHandler{
		orderService:     orderService,
		notificationPool: pool,
	}
}

func (h *OrderHandler) CreateOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("USER ID", claims.UserID)

	orderID, err := h.orderService.CreateOrder(
		r.Context(),
		claims.UserID,
	)

	h.notificationPool.Submit(
		service.NotificationJob{
			OrderID: orderID,
			UserID:  claims.UserID,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		map[string]int64{
			"order_id": orderID,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetOrders(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrders(
		r.Context(),
		claims.UserID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		map[string][]model.OrderInfo{
			"orders": orders,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetOrderByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orderIDStr := r.PathValue("id")

	orderID, err := strconv.ParseInt(
		orderIDStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "Inalid ID", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(
		r.Context(),
		claims.UserID,
		orderID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (h *OrderHandler) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPatch {
		http.Error(w, "only PATCH method", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orderStr := r.PathValue("id")

	orderID, err := strconv.ParseInt(
		orderStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.orderService.UpdateStatus(
		r.Context(),
		orderID,
		req.Status,
		claims.UserID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		map[string]string{
			"message": "status updated",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}
}
