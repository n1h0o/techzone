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
