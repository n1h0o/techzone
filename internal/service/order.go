package service

import (
	"context"
	"errors"
	"log"
	"techzone/internal/event"
	"techzone/internal/metrics"
	"techzone/internal/model"
	"techzone/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	publisher   EventPublisher
	db          *pgxpool.Pool
	redis       *redis.Client
}

func NewOrderService(
	orderRepo OrderRepository,
	cartRepo CartRepository,
	productRepo ProductRepository,
	publisher EventPublisher,
	db *pgxpool.Pool,
	redis *redis.Client,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		publisher:   publisher,
		db:          db,
		redis:       redis,
	}
}

type OrderDetails struct {
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
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback skipped: %v", err)
		}
	}()

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

	log.Println("commit transaction")

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	log.Println("transaction committed")

	if err := s.redis.Del(ctx, "products").Err(); err != nil {
		log.Printf("failed to clear products cache: %v", err)
	}

	log.Println("publishing order.created")

	err = s.publisher.PublishOrderCreated(
		ctx,
		event.OrderCreatedEvent{
			OrderID: orderID,
			UserID:  userID,
		},
	)

	log.Println("publish finished")
	if err != nil {
		log.Printf(
			"failed to publish order.created event: %v",
			err,
		)
	}
	log.Printf(
		"order created id=%d user=%d total=%.2f",
		orderID,
		userID,
		sum,
	)
	metrics.OrdersCreatedTotal.Inc()
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
) (*OrderDetails, error) {
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

	return &OrderDetails{
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

	case model.OrderNew:
		if status != model.OrderProcessing {
			return errors.New("invalid status transition")
		}

	case model.OrderProcessing:
		if status != model.OrderCompleted {
			return errors.New("invalid status transition")
		}

	case model.OrderCompleted:
		return errors.New("order already completed")

	case model.OrderCancelled:
		return errors.New("order is cancelled")

	default:
		return errors.New("unknown order status")
	}

	return s.orderRepo.UpdateStatus(
		ctx,
		orderID,
		status,
	)
}
