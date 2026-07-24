package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"techzone/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
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

	UpdateProduct(
		ctx context.Context,
		product *model.Product,
	) error

	SetActive(
		ctx context.Context,
		productID int64,
		active bool,
	) error

	GetAllForAdmin(
		ctx context.Context,
	) ([]model.Product, error)
}

// инкапсулирует правила каталога и слой кэширования
type ProductService struct {
	productRepo ProductRepository
	redis       *redis.Client
}

// создает сервис каталога
func NewProductService(
	productRepo ProductRepository,
	redis *redis.Client,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		redis:       redis,
	}
}

// валидирует данные товара и сбрасывает кэш после записи
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
		ImageURL:    input.ImageURL,
	}

	log.Printf(
		"product created id=%d name=%s",
		product.ID,
		product.Name,
	)
	productID, err := s.productRepo.Create(
		ctx,
		product,
	)
	if err != nil {
		return 0, err
	}

	_ = s.redis.Del(ctx, "products").Err()

	return productID, nil
}

// возвращает один активный товар
func (s *ProductService) GetProduct(
	ctx context.Context,
	id int64,
) (*model.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// сначала читает каталог из кэша и только потом идет в postgres
func (s *ProductService) GetProducts(
	ctx context.Context,
) ([]model.Product, error) {

	data, err := s.redis.Get(ctx, "products").Bytes()

	if err == nil {
		var products []model.Product

		if err := json.Unmarshal(data, &products); err == nil {
			log.Println("products loaded from redis")
			return products, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("redis error: %v", err)
	}

	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(products)
	if err == nil {
		// короткий ttl позволяет пережить редкие промахи инвалидции
		_ = s.redis.Set(
			ctx,
			"products",
			bytes,
			5*time.Minute,
		).Err()
	}
	log.Println("products loaded from postgres")

	return products, nil
}

// отдает полный каталог для админки без фильтра по активности
func (s *ProductService) GetProductsForAdmin(
	ctx context.Context,
) ([]model.Product, error) {
	return s.productRepo.GetAllForAdmin(ctx)
}

// обновляет товар и инвалидирует общий список каталога
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	productID int64,
	input CreateProductInput,
) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("empty name")
	}
	if len(input.Name) > 255 {
		return errors.New("name too long")
	}
	if strings.TrimSpace(input.Description) == "" {
		return errors.New("empty description")
	}
	if input.Price <= 0 {
		return errors.New("invalid price")
	}
	if input.Stock < 0 {
		return errors.New("invalid stock")
	}

	product := &model.Product{
		ID:          productID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		ImageURL:    input.ImageURL,
	}

	log.Printf("product updated id = %d", productID)

	err := s.productRepo.UpdateProduct(
		ctx,
		product,
	)

	if err != nil {
		return err
	}

	_ = s.redis.Del(ctx, "products").Err()

	return nil
}

// меняет флаг активности и сбрасывает кэш каталога
func (s *ProductService) SetProductStatus(
	ctx context.Context,
	productID int64,
	active bool,
) error {

	if productID <= 0 {
		return errors.New("invalid product id")
	}

	log.Printf("product deleted id=%d", productID)

	err := s.productRepo.SetActive(
		ctx,
		productID,
		active,
	)
	if err != nil {
		return err
	}

	_ = s.redis.Del(ctx, "products").Err()

	return nil
}
