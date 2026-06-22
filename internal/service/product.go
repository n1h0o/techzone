package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"techzone/internal/model"
)

type ProductRepository interface {
	Create(
		ctx context.Context,
		product *model.Product,
	) (int64, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*model.Product, error)

	GetAll(
		ctx context.Context,
	) ([]model.Product, error)

	DecreaseStock(
		ctx context.Context,
		productID int64,
		quantity int,
	) error
}

type ProductService struct {
	productRepo ProductRepository
}

func NewProductService(
	productRepo ProductRepository,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(
	ctx context.Context,
	input CreateProductInput,
) (int64, error) {
	if strings.TrimSpace(input.Name) == "" {
		return 0, errors.New("empty name")
	}
	if len(input.Name) > 255 {
		return 0, errors.New("name too long")
	}
	if strings.TrimSpace(input.Description) == "" {
		return 0, errors.New("empty description")
	}
	if input.Price <= 0 {
		return 0, errors.New("invalid price")
	}
	if input.Stock < 0 {
		return 0, errors.New("invalid stock")
	}

	product := &model.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	log.Printf(
		"product created id=%d name=%s",
		product.ID,
		product.Name,
	)
	return s.productRepo.Create(
		ctx,
		product,
	)
}

func (s *ProductService) GetProduct(
	ctx context.Context,
	id int64,
) (*model.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) GetProducts(
	ctx context.Context,
) ([]model.Product, error) {
	return s.productRepo.GetAll(ctx)
}
