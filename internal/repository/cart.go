package repository

import (
	"context"
	"errors"
	"techzone/internal/model"
)

type CartRepository struct {
	db DBTX
}

func NewCartRepository(
	db DBTX,
) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) (*model.Cart, error) {

	var cart model.Cart

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, user_id, created_at
		FROM carts
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) Create(
	ctx context.Context,
	userID int64,
) (int64, error) {

	var id int64

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO carts(user_id)
		VALUES ($1)
		RETURNING id
		`,
		userID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CartRepository) AddItem(
	ctx context.Context,
	cartID int64,
	productID int64,
	quantity int,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO cart_items(
		cart_id,
		product_id,
		quantity
		)
		VALUES($1,$2,$3)
		`,
		cartID,
		productID,
		quantity,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) GetCart(
	ctx context.Context,
	cartID int64,
) ([]model.CartItemInfo, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT 
		ci.id,
		p.id,
		p.name,
		p.price,
		ci.quantity
		FROM cart_items ci
		JOIN products p
			ON p.id = ci.product_id
		WHERE ci.cart_id = $1
		`,
		cartID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cart []model.CartItemInfo

	for rows.Next() {
		var item model.CartItemInfo
		if err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Name,
			&item.Price,
			&item.Quantity,
		); err != nil {
			return nil, err
		}
		cart = append(cart, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *CartRepository) DeleteItem(
	ctx context.Context,
	itemID int64,
	cartID int64,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		DELETE FROM cart_items
		WHERE id = $1
		AND cart_id = $2
		`,
		itemID, cartID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *CartRepository) ClearCart(
	ctx context.Context,
	cartID int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM cart_items
		WHERE cart_id = $1
		`,
		cartID,
	)
	if err != nil {
		return err
	}
	return nil
}
