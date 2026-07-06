package repository

import (
	"context"
	"errors"
	"techzone/internal/model"
)

type OrderRepository struct {
	db DBTX
}

func NewOrderRepository(
	db DBTX,
) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	order *model.Order,
) (int64, error) {

	var id int64

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO orders(
		user_id,
		total_price
		)
		VALUES ($1,$2)
		RETURNING id
		`,
		order.UserID,
		order.TotalPrice,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrderRepository) CreateItem(
	ctx context.Context,
	item *model.OrderItem,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO order_items(
		order_id,
		product_id,
		quantity,
		price
		)
		VALUES($1,$2,$3,$4)
		`,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]model.OrderInfo, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			o.id,
			o.status,
			COALESCE(p.status,'not_paid') AS payment_status, 
			o.total_price
		FROM orders o
		LEFT JOIN payments p
			ON p.order_id = o.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderInfo

	for rows.Next() {
		var order model.OrderInfo
		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.PaymentStatus,
			&order.TotalPrice,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetByID(
	ctx context.Context,
	orderID int64,
	userID int64,
) (*model.Order, error) {

	var order model.Order

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
		id,
		user_id,
		status,
		total_price,
		created_at
		FROM orders
		WHERE id = $1
		AND user_id = $2
		`,
		orderID,
		userID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetItems(
	ctx context.Context,
	orderID int64,
) ([]model.OrderItemInfo, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
		p.id,
		p.name,
		oi.price,
		oi.quantity
		FROM order_items oi
		JOIN products p
		ON p.id = oi.product_id
		WHERE oi.order_id = $1
		`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItemInfo

	for rows.Next() {
		var item model.OrderItemInfo
		if err := rows.Scan(
			&item.ProductID,
			&item.Name,
			&item.Price,
			&item.Quantity,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	orderID int64,
	status string,
) error {

	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE orders
		SET status = $1
		WHERE id = $2
		`,
		status,
		orderID,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *OrderRepository) LockOrder(
	ctx context.Context,
	orderID int64,
) error {
	var id int64

	return r.db.QueryRow(
		ctx,
		`
		SELECT id
		FROM orders
		WHERE id = $1
		FOR UPDATE
		`,
		orderID,
	).Scan(&id)
}
