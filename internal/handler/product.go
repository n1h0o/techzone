package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"techzone/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(
	productService *service.ProductService,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts godoc
//
// @Summary Получить список товаров
// @Description Возвращает список активных товаров
// @Tags products
// @Produce json
// @Success 200 {array} model.Product
// @Router /products [get]
func (h *ProductHandler) GetProducts(
	w http.ResponseWriter, r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}
	products, err := h.productService.GetProducts(
		r.Context(),
	)
	if err != nil {
		log.Printf("GetProducts failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// GetProductsForAdmin godoc
//
// @Summary Получить список всех товаров
// @Description Возвращает полный список товаров, включая неактивные. Доступно только администраторам.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Product
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /admin/products [get]
func (h *ProductHandler) GetProductsForAdmin(
	w http.ResponseWriter, r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}
	products, err := h.productService.GetProductsForAdmin(
		r.Context(),
	)
	if err != nil {
		log.Printf("GetProductsForAdmin failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateProduct godoc
//
// @Summary Создать товар
// @Description Создает новый товар
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateProductInput true "Данные товара"
// @Success 201 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method", http.StatusMethodNotAllowed)
		return
	}

	var prod service.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	productID, err := h.productService.CreateProduct(
		r.Context(),
		prod,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(
		ProductResponse{
			ID: productID,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// GetProduct godoc
//
// @Summary Получить товар
// @Description Возвращает товар по ID
// @Tags products
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} model.Product
// @Failure 404 {string} string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProduct(
		r.Context(),
		id,
	)

	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

// UpdateProduct godoc
//
// @Summary Обновить товар
// @Description Полностью обновляет товар
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param request body service.CreateProductInput true "Новые данные"
// @Success 200 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		http.Error(w, "only PUT method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var prod service.CreateProductInput

	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.productService.UpdateProduct(
		r.Context(),
		id,
		prod,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		MessageResponse{
			Message: "product was updated",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// SetProductStatus godoc
//
// @Summary Изменить статус товара
// @Description Активирует или деактивирует товар
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param request body handler.ProductStatusInput true "Статус товара"
// @Success 200 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Router /products/{id}/status [patch]
func (h *ProductHandler) SetProductStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPatch {
		http.Error(w, "only PATCH method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req ProductStatusInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.productService.SetProductStatus(
		r.Context(),
		id,
		req.IsActive,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		MessageResponse{
			Message: "status updated",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
