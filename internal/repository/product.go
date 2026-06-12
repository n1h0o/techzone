package repository

import (
	"context"
	"errors"
	"log"
	"techzone/internal/model"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(
	db DBTX,
) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(
	ctx context.Context,
	product *model.Product,
) (int64, error) {
	var id int64

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO products(
			name,
			description,
			price,
			stock
		)
		VALUES ($1,$2,$3,$4)
		RETURNING id
		`,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	).Scan(&id)

	return id, err
}

func (r *ProductRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.Product, error) {

	var product model.Product

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, name, description, price, stock,created_at
		FROM products
		WHERE id = $1
		`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(
	ctx context.Context,
) ([]model.Product, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT id, name, description, price, stock, created_at
		FROM products
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var prod model.Product
		if err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.Price,
			&prod.Stock,
			&prod.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) DecreaseStock(
	ctx context.Context,
	productID int64,
	quantity int,
) error {

	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
		AND stock >= $1
		`,
		quantity, productID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		log.Printf(
			"not enough stock product=%d quantity=%d",
			productID,
			quantity,
		)
		return errors.New("insufficient stock")
	}
	return nil
}
