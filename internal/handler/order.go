package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"techzone/internal/middleware"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(
	orderService *service.OrderService,
) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
//
// @Summary Создать заказ
// @Description Создает заказ из текущей корзины пользователя
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.OrderCreateResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /orders [post]
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		OrderCreateResponse{
			OrderID: orderID,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// GetOrders godoc
//
// @Summary Получить список заказов
// @Description Возвращает все заказы текущего пользователя
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.OrdersResponse
// @Failure 401 {string} string
// @Router /orders [get]
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
		OrdersResponse{
			Orders: orders,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// GetOrderByID godoc
//
// @Summary Получить заказ
// @Description Возвращает заказ вместе с товарами
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} service.OrderDetails
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Router /orders/{id} [get]
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

// UpdateStatus godoc
//
// @Summary Обновить статус заказа
// @Description Изменяет статус заказа (только для администратора)
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Param request body handler.UpdateStatusRequest true "Новый статус"
// @Success 200 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Router /orders/{id}/status [patch]
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
		MessageResponse{
			Message: "status updated",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}
}
