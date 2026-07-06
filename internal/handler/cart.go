package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techzone/internal/middleware"
	"techzone/internal/model"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(
	cartService *service.CartService,
) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddToCart godoc
//
// @Summary Добавить товар в корзину
// @Description Добавляет товар в корзину пользователя
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.AddToCartInput true "Товар"
// @Success 200 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method", http.StatusMethodNotAllowed)
		return
	}

	var req service.AddToCartInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)

	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err := h.cartService.AddToCart(
		r.Context(),
		claims.UserID,
		req.ProductID,
		req.Quantity,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		MessageResponse{
			Message: "product added to cart",
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// GetCart godoc
//
// @Summary Получить корзину
// @Description Возвращает корзину текущего пользователя
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.CartResponse
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /cart [get]
func (h *CartHandler) GetCart(
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

	cart, err := h.cartService.GetCart(r.Context(), claims.UserID)

	if cart == nil {
		cart = []model.CartItemInfo{}
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		CartResponse{
			Items: cart,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// DeleteItem godoc
//
// @Summary Удалить товар из корзины
// @Description Удаляет товар из корзины пользователя
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Param item_id path int true "ID элемента корзины"
// @Success 200 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Router /cart/items/{item_id} [delete]
func (h *CartHandler) DeleteItem(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodDelete {
		http.Error(w, "only DELETE method", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	itemIDStr := r.PathValue("item_id")

	itemID, err := strconv.ParseInt(
		itemIDStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.cartService.DeleteItem(
		r.Context(),
		claims.UserID,
		itemID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		MessageResponse{
			Message: "item deleted",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
