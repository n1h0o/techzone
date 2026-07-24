package repository

import (
	"context"
	"errors"
	"log"
	"techzone/internal/model"
	"techzone/pkg/dbtx"
)

type ProductRepository struct {
	db dbtx.DBTX
}

func NewProductRepository(
	db dbtx.DBTX,
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
			stock,
			image_url
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
		`,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
	).Scan(&id)

	if err != nil {
		log.Printf("create product error: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *ProductRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.Product, error) {

	var product model.Product

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, name, description, price, stock,created_at, image_url
		FROM products
		WHERE id = $1
		AND is_active = true
		LIMIT 1
		`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.ImageURL,
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
		SELECT id, name, description, price, stock, created_at, image_url
		FROM products
		WHERE is_active = true
		ORDER BY id DESC;
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
			&prod.ImageURL,
		); err != nil {
			log.Printf("SCAN ERROR: %v", err)
			return nil, err
		}
		products = append(products, prod)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ROWS ERROR: %v", err)
		return nil, err
	}
	if products == nil {
		products = []model.Product{}
	}
	return products, nil
}

func (r *ProductRepository) GetAllForAdmin(
	ctx context.Context,
) ([]model.Product, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT id, name, description, price, stock, created_at, image_url, is_active
		FROM products
		ORDER BY id DESC;
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
			&prod.ImageURL,
			&prod.IsActive,
		); err != nil {
			log.Printf("SCAN ERROR: %v", err)
			return nil, err
		}
		products = append(products, prod)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ROWS ERROR: %v", err)
		return nil, err
	}

	if products == nil {
		products = []model.Product{}
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
		AND is_active = true
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

func (r *ProductRepository) UpdateProduct(
	ctx context.Context,
	product *model.Product,
) error {

	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE products
		SET
			name = $2,
			description = $3,
			price = $4,
			stock = $5,
			image_url = $6
		WHERE id = $1
		`,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *ProductRepository) SetActive(
	ctx context.Context,
	productID int64,
	active bool,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE products
		SET is_active = $2
		WHERE id = $1
		`,
		productID,
		active,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
