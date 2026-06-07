package repository

import (
	"context"
	"techzone/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(
	db *pgxpool.Pool,
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
