package handler

import (
	"encoding/json"
	"net/http"
	"techzone/internal/middleware"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(
	paymentService *service.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Pay godoc
//
// @Summary Оплатить заказ
// @Description Выполняет оплату заказа с использованием Idempotency-Key
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body handler.PaymentRequest true "Данные платежа"
// @Param Idempotency-Key header string true "Уникальный ключ идемпотентности"
// @Success 200 {object} model.Payment
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /payments [post]
func (h *PaymentHandler) Pay(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req PaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header is empty", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.Pay(
		r.Context(),
		claims.UserID,
		req.OrderID,
		idempotencyKey,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
