package service

import (
	"context"
	"errors"
	"log"
	"techzone/internal/model"
	"techzone/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(
		ctx context.Context,
		order *model.Order,
	) (int64, error)

	CreateItem(
		ctx context.Context,
		item *model.OrderItem,
	) error

	GetByUserID(
		ctx context.Context,
		userID int64,
	) ([]model.OrderInfo, error)

	GetByID(
		ctx context.Context,
		orderID int64,
		userID int64,
	) (*model.Order, error)

	GetItems(
		ctx context.Context,
		orderID int64,
	) ([]model.OrderItemInfo, error)

	UpdateStatus(
		ctx context.Context,
		orderID int64,
		status string,
	) error
}

type OrderService struct {
	orderRepo   OrderRepository
	cartRepo    CartRepository
	productRepo ProductRepository
	db          *pgxpool.Pool
}

func NewOrderService(
	orderRepo OrderRepository,
	cartRepo CartRepository,
	productRepo ProductRepository,
	db *pgxpool.Pool,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		db:          db,
	}
}

type orderDetails struct {
	Order *model.Order          `json:"order"`
	Items []model.OrderItemInfo `json:"items"`
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userID int64,
) (int64, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Printf(
			"order creation failed user=%d err=%v",
			userID,
			err,
		)
		return 0, err
	}
	orderRepo := repository.NewOrderRepository(tx)
	cartRepo := repository.NewCartRepository(tx)
	productRepo := repository.NewProductRepository(tx)
	defer tx.Rollback(ctx)

	var sum float64

	cartInfo, err := cartRepo.GetByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return 0, err
	}
	items, err := cartRepo.GetCart(
		ctx,
		cartInfo.ID,
	)
	if err != nil {
		return 0, err
	}
	if len(items) == 0 {
		return 0, errors.New("cart is empty")
	}

	for _, item := range items {
		sum += item.Price * float64(item.Quantity)
	}
	orderID, err := orderRepo.Create(
		ctx,
		&model.Order{
			UserID:     userID,
			TotalPrice: sum,
		},
	)
	if err != nil {
		return 0, err
	}
	for _, item := range items {
		err := orderRepo.CreateItem(
			ctx,
			&model.OrderItem{
				OrderID:   orderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			},
		)
		if err != nil {
			return 0, err
		}
		err = productRepo.DecreaseStock(
			ctx,
			item.ProductID,
			item.Quantity,
		)
		if err != nil {
			return 0, err
		}
	}

	err = cartRepo.ClearCart(
		ctx,
		cartInfo.ID,
	)

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	log.Printf(
		"order created id=%d user=%d total=%.2f",
		orderID,
		userID,
		sum,
	)
	return orderID, nil
}

func (s *OrderService) GetOrders(
	ctx context.Context,
	userID int64,
) ([]model.OrderInfo, error) {
	return s.orderRepo.GetByUserID(
		ctx, userID,
	)
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*orderDetails, error) {
	order, err := s.orderRepo.GetByID(
		ctx,
		orderID,
		userID,
	)
	if err != nil {
		return nil, err
	}
	items, err := s.orderRepo.GetItems(
		ctx,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	return &orderDetails{
		Order: order,
		Items: items,
	}, nil
}

func (s *OrderService) UpdateStatus(
	ctx context.Context,
	orderID int64,
	status string,
	userID int64,
) error {

	order, err := s.orderRepo.GetByID(
		ctx,
		orderID,
		userID,
	)
	if err != nil {
		return err
	}

	switch order.Status {
	case "new":
		if status != "processing" {
			return errors.New("invalid status transition")
		}
	case "processing":
		if status != "completed" {
			return errors.New("invalid status transition")
		}
	case "completed":
		return errors.New("order already completed")
	default:
		return errors.New("invalid status")
	}
	return s.orderRepo.UpdateStatus(
		ctx,
		orderID,
		status,
	)
}
