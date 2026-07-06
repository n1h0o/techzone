package service

import (
	"context"
	"errors"
	"log"
	"techzone/internal/event"
	"techzone/internal/model"
	"techzone/internal/payment"
	"techzone/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(
		ctx context.Context,
		payment *model.Payment,
	) (int64, error)

	GetByIdempotencyKey(
		ctx context.Context,
		key string,
	) (*model.Payment, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*model.Payment, error)

	UpdateStatus(
		ctx context.Context,
		paymentID int64,
		status string,
		transactionID string,
	) error

	GetByOrderID(
		ctx context.Context,
		orderID int64,
	) (*model.Payment, error)
}

type OrderProvider interface {
	GetByID(
		ctx context.Context,
		orderID int64,
		userID int64,
	) (*model.Order, error)
}

type PaymentService struct {
	gateway   payment.Gateway
	publisher EventPublisher

	db *pgxpool.Pool
}

func NewPaymentService(
	gateway payment.Gateway,
	publisher EventPublisher,
	db *pgxpool.Pool,
) *PaymentService {
	return &PaymentService{
		gateway:   gateway,
		publisher: publisher,
		db:        db,
	}
}

func (s *PaymentService) findExistingPayment(
	ctx context.Context,
	repo PaymentRepository,
	orderID int64,
	key string,
) (*model.Payment, error) {
	payment, err := repo.GetByIdempotencyKey(
		ctx,
		key,
	)
	if err == nil {
		return payment, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return repo.GetByOrderID(
		ctx,
		orderID,
	)
}

func (s *PaymentService) createPayment(
	ctx context.Context,
	repo PaymentRepository,
	order *model.Order,
	userID int64,
	key string,
) (*model.Payment, error) {

	payment := &model.Payment{
		OrderID:        order.ID,
		UserID:         userID,
		Amount:         int64(order.TotalPrice),
		Status:         model.PaymentPending,
		IdempotencyKey: key,
	}

	id, err := repo.Create(
		ctx,
		payment,
	)

	if err != nil {
		return nil, err
	}

	payment.ID = id

	return payment, nil
}

func (s *PaymentService) processGateway(
	ctx context.Context,
	repo PaymentRepository,
	pay *model.Payment,
) error {
	resp, err := s.gateway.Pay(
		ctx,
		payment.PaymentRequest{
			OrderID: pay.OrderID,
			UserID:  pay.UserID,
			Amount:  pay.Amount,
		},
	)
	if err != nil {
		if updateErr := repo.UpdateStatus(
			ctx,
			pay.ID,
			model.PaymentFailed,
			"",
		); updateErr != nil {
			return updateErr
		}
		pay.Status = model.PaymentFailed

		return err
	}

	pay.TransactionID = resp.TransactionID

	if resp.Status == payment.StatusSuccess {

		if err := repo.UpdateStatus(
			ctx,
			pay.ID,
			model.PaymentSuccess,
			resp.TransactionID,
		); err != nil {
			return err
		}

		pay.Status = model.PaymentSuccess

		return nil
	}

	if err := repo.UpdateStatus(
		ctx,
		pay.ID,
		model.PaymentFailed,
		resp.TransactionID,
	); err != nil {
		return err
	}

	pay.Status = model.PaymentFailed

	return err
}

func (s *PaymentService) Pay(
	ctx context.Context,
	userID int64,
	orderID int64,
	idempotencyKey string,
) (*model.Payment, error) {

	if idempotencyKey == "" {
		return nil, errors.New("idempotency key is empty")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback skipped: %v", err)
		}
	}()

	paymentRepo := repository.NewPaymentRepository(tx)
	orderRepo := repository.NewOrderRepository(tx)

	if err := orderRepo.LockOrder(ctx, orderID); err != nil {
		return nil, err
	}

	pay, err := s.findExistingPayment(
		ctx,
		paymentRepo,
		orderID,
		idempotencyKey,
	)

	switch {
	case err == nil:
		return pay, nil

	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	order, err := orderRepo.GetByID(
		ctx,
		orderID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	pay, err = s.createPayment(
		ctx,
		paymentRepo,
		order,
		userID,
		idempotencyKey,
	)
	if err != nil {
		return nil, err
	}

	if err := s.processGateway(
		ctx,
		paymentRepo,
		pay,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	if pay.Status == model.PaymentSuccess {
		if err := s.publisher.PublishPaymentCompleted(
			ctx,
			event.PaymentCompletedEvent{
				PaymentID: pay.ID,
				OrderID:   pay.OrderID,
				UserID:    pay.UserID,
			},
		); err != nil {
			log.Printf("failed to publish payment.completed: %v", err)
		}
	}
	return pay, nil
}
