package service

import (
	"context"
	"errors"
	"techzone/internal/model"
	"testing"
)

type MockOrderRepo struct {
	Orders    []model.OrderInfo
	OrdersErr error

	Order    *model.Order
	OrderErr error

	UpdateCalled bool
	UpdateErr    error

	Items    []model.OrderItemInfo
	ItemsErr error
}

func (m *MockOrderRepo) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]model.OrderInfo, error) {
	return m.Orders, m.OrdersErr
}

func (m *MockOrderRepo) GetByID(
	ctx context.Context,
	orderID int64,
	userID int64,
) (*model.Order, error) {
	return m.Order, m.OrderErr
}

func (m *MockOrderRepo) GetItems(
	ctx context.Context,
	orderID int64,
) ([]model.OrderItemInfo, error) {
	return m.Items, m.ItemsErr
}

func (m *MockOrderRepo) UpdateStatus(
	ctx context.Context,
	orderID int64,
	status string,
) error {
	m.UpdateCalled = true
	return m.UpdateErr
}

func (m *MockOrderRepo) Create(
	ctx context.Context,
	order *model.Order,
) (int64, error) {
	return 1, nil
}

func (m *MockOrderRepo) CreateItem(
	ctx context.Context,
	item *model.OrderItem,
) error {
	return nil
}

func TestUpdateStatus_NewToProcessing(t *testing.T) {

	repo := &MockOrderRepo{
		Order: &model.Order{
			ID:     1,
			Status: "new",
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	err := service.UpdateStatus(
		context.Background(),
		1,
		"processing",
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.UpdateCalled {
		t.Fatal("UpdateStatus was not called")
	}

}

func TestUpdateStatus_NewToCompleted(t *testing.T) {

	repo := &MockOrderRepo{
		Order: &model.Order{
			ID:     1,
			Status: "new",
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	err := service.UpdateStatus(
		context.Background(),
		1,
		"completed",
		1,
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if repo.UpdateCalled {
		t.Fatal("UpdateStatus was not called")
	}

}

func TestUpdateStatus_ProcessingToCompleted(t *testing.T) {

	repo := &MockOrderRepo{
		Order: &model.Order{
			ID:     1,
			Status: "processing",
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	err := service.UpdateStatus(
		context.Background(),
		1,
		"completed",
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.UpdateCalled {
		t.Fatal("UpdateStatus was not called")
	}
}

func TestUpdateStatus_CompletedOrder(t *testing.T) {
	repo := &MockOrderRepo{
		Order: &model.Order{
			ID:     1,
			Status: "completed",
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	err := service.UpdateStatus(
		context.Background(),
		1,
		"Processing",
		1,
	)
	if err == nil {
		t.Fatal(err)
	}

	if repo.UpdateCalled {
		t.Fatal("Update was not called")
	}
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	repo := &MockOrderRepo{
		Order: &model.Order{
			ID:     1,
			Status: "unknown",
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	err := service.UpdateStatus(
		context.Background(),
		1,
		"ok",
		1,
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if repo.UpdateCalled {
		t.Fatal("Update was not called")
	}
}

func TestGetOrders_Success(t *testing.T) {
	repo := &MockOrderRepo{
		Orders: []model.OrderInfo{
			{
				ID: 1,
			},
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	orders, err := service.GetOrders(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatal("expected error")
	}

	if len(orders) != 1 {
		t.Fatal("expected 1 order")
	}

}

func TestGetOrders_Error(t *testing.T) {
	repo := &MockOrderRepo{
		OrdersErr: errors.New("db error"),
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	_, err := service.GetOrders(
		context.Background(),
		1,
	)

	if err == nil {
		t.Fatal("expected error")
	}

}

func TestGetOrder_Success(t *testing.T) {
	repo := &MockOrderRepo{
		Order: &model.Order{
			ID: 1,
		},
		Items: []model.OrderItemInfo{
			{
				ProductID: 1,
			},
		},
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	result, err := service.GetOrder(
		context.Background(),
		1,
		1,
	)

	if err != nil {
		t.Fatal("expected error")
	}

	if result == nil {
		t.Fatal("expected order")
	}

	if len(result.Items) != 1 {
		t.Fatal("expected 1 item")
	}

}

func TestGetOrder_GetBByIDError(t *testing.T) {
	repo := &MockOrderRepo{
		OrderErr: errors.New("db error"),
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	_, err := service.GetOrder(
		context.Background(),
		1,
		1,
	)

	if err == nil {
		t.Fatal("expected error")
	}

}

func TestGetOrders_GetItemsError(t *testing.T) {
	repo := &MockOrderRepo{
		Order: &model.Order{
			ID: 1,
		},
		ItemsErr: errors.New("db error"),
	}

	service := NewOrderService(
		repo,
		nil,
		nil,
		nil,
	)

	_, err := service.GetOrder(
		context.Background(),
		1,
		1,
	)

	if err == nil {
		t.Fatal("expected error")
	}

}
