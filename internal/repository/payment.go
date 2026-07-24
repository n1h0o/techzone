package repository

import (
	"context"
	"errors"
	"techzone/internal/model"
	"techzone/pkg/dbtx"
)

type PaymentRepository struct {
	db dbtx.DBTX
}

func NewPaymentRepository(
	db dbtx.DBTX,
) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(
	ctx context.Context,
	payment *model.Payment,
) (int64, error) {

	var id int64

	err := r.db.QueryRow(
		ctx,

		`
		INSERT INTO payments(
			order_id,
			user_id,
			amount,
			status,
			idempotency_key,
			transaction_id
		)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id
		`,
		payment.OrderID,
		payment.UserID,
		payment.Amount,
		payment.Status,
		payment.IdempotencyKey,
		payment.TransactionID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil

}

func (r *PaymentRepository) GetByIdempotencyKey(
	ctx context.Context,
	key string,
) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			order_id,
			user_id,
			amount,
			status,
			idempotency_key,
			transaction_id,
			created_at
			FROM payments
		WHERE idempotency_key = $1
		`,
		key,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.IdempotencyKey,
		&payment.TransactionID,
		&payment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil

}

func (r *PaymentRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.Payment, error) {
	var payment model.Payment

	err := r.db.QueryRow(
		ctx,
		`
			SELECT 
				id, 
				order_id,
				user_id,
				amount,
				status,
				idempotency_key,
				transaction_id,
				created_at
				FROM payments
			WHERE id = $1
			`,
		id,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.IdempotencyKey,
		&payment.TransactionID,
		&payment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status string,
	transactionID string,
) error {

	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET 
			status = $1,
			transaction_id = $2

		WHERE id = $3
		`,
		status,
		transactionID,
		id,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (r *PaymentRepository) GetByOrderID(
	ctx context.Context,
	orderID int64,
) (*model.Payment, error) {
	var payment model.Payment

	err := r.db.QueryRow(
		ctx,
		`
			SELECT 
				id, 
				order_id,
				user_id,
				amount,
				status,
				idempotency_key,
				transaction_id,
				created_at
				FROM payments
			WHERE order_id = $1
			AND status = 'success'
			LIMIT 1
			`,
		orderID,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.IdempotencyKey,
		&payment.TransactionID,
		&payment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
