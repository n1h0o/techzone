package service

import (
	"context"
	"errors"
	"techzone/internal/model"
	"testing"

	"github.com/jackc/pgx/v5"
)

type MockCartRepo struct {
	Cart    *model.Cart
	CartErr error

	CreatedID    int64
	CreateErr    error
	CreateCalled bool

	AddErr        error
	AddItemCalled bool

	Items      []model.CartItemInfo
	GetCartErr error

	DeleteErr    error
	DeleteCalled bool
}

func (m *MockCartRepo) GetByUserID(
	ctx context.Context,
	userID int64,
) (*model.Cart, error) {
	return m.Cart, m.CartErr
}

func (m *MockCartRepo) Create(
	ctx context.Context,
	userID int64,
) (int64, error) {
	m.CreateCalled = true
	return m.CreatedID, m.CreateErr
}

func (m *MockCartRepo) AddItem(
	ctx context.Context,
	cartID int64,
	productID int64,
	quantity int,
) error {
	m.AddItemCalled = true
	return m.AddErr
}

func (m *MockCartRepo) GetCart(
	ctx context.Context,
	cartID int64,
) ([]model.CartItemInfo, error) {
	return m.Items, m.GetCartErr
}

func (m *MockCartRepo) DeleteItem(
	ctx context.Context,
	itemID int64,
	cartID int64,
) error {
	m.DeleteCalled = true
	return m.DeleteErr
}

func TestAddToCart_InvalidProductID(
	t *testing.T,
) {
	repo := &MockCartRepo{}

	service := NewCartService(repo)

	err := service.AddToCart(
		context.Background(),
		1,
		0,
		1,
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid product id" {
		t.Fatalf(
			"expected invalid product id, got %v",
			err,
		)
	}
}

func TestAddToCart_InvalidQuantity(
	t *testing.T,
) {
	repo := &MockCartRepo{}

	service := NewCartService(repo)

	err := service.AddToCart(
		context.Background(),
		1,
		1,
		0,
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "invalid quantity" {
		t.Fatalf(
			"expected invalid quantity, got %v",
			err,
		)
	}
}

func TestAddToCart_CreateCartIfNotExists(
	t *testing.T,
) {
	repo := &MockCartRepo{
		CartErr:   pgx.ErrNoRows,
		CreatedID: 10,
	}

	service := NewCartService(repo)

	err := service.AddToCart(
		context.Background(),
		1,
		1,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.CreateCalled {
		t.Fatal("Create was not called")
	}

	if !repo.AddItemCalled {
		t.Fatal("AddItem was not called")
	}

}

func TestAddToCart_ExistingCart(
	t *testing.T,
) {
	repo := &MockCartRepo{
		Cart: &model.Cart{
			ID: 5,
		},
	}

	service := NewCartService(repo)

	err := service.AddToCart(
		context.Background(),
		1,
		1,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.AddItemCalled {
		t.Fatal("AddItem was not called")
	}
}

func TestGetCart_EmptyCart(
	t *testing.T,
) {
	repo := &MockCartRepo{
		CartErr: pgx.ErrNoRows,
	}

	service := NewCartService(repo)

	items, err := service.GetCart(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 0 {
		t.Fatal("expected empty cart")
	}
}

func TestGetCart_RepositoryError(
	t *testing.T,
) {
	repo := &MockCartRepo{
		CartErr: errors.New("db error"),
	}

	service := NewCartService(repo)

	_, err := service.GetCart(
		context.Background(),
		1,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteItem(
	t *testing.T,
) {
	repo := &MockCartRepo{
		Cart: &model.Cart{
			ID: 1,
		},
	}

	service := NewCartService(repo)

	err := service.DeleteItem(
		context.Background(),
		1,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !repo.DeleteCalled {
		t.Fatal("DeleteItem was not called")
	}
}
