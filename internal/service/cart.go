package service

import (
	"context"
	"errors"
	"log"
	"techzone/internal/model"

	"github.com/jackc/pgx/v5"
)

type CartRepository interface {
	AddItem(
		ctx context.Context,
		CartID int64,
		ProductID int64,
		Quantity int,
	) error
	GetCart(
		ctx context.Context,
		cartID int64,
	) ([]model.CartItemInfo, error)

	GetByUserID(
		ctx context.Context,
		userID int64,
	) (*model.Cart, error)

	Create(
		ctx context.Context,
		userID int64,
	) (int64, error)

	DeleteItem(
		ctx context.Context,
		itemID int64,
		cartID int64,
	) error
}

type CartService struct {
	cartRepo CartRepository
}

func NewCartService(
	cartRepo CartRepository,
) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

func (s *CartService) AddToCart(
	ctx context.Context,
	userID int64,
	productID int64,
	quantity int,
) error {
	if productID <= 0 {
		return errors.New("invalid product id")
	}
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	cart, err := s.cartRepo.GetByUserID(
		ctx,
		userID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		cartID, err := s.cartRepo.Create(
			ctx,
			userID,
		)
		if err != nil {
			return err
		}
		return s.cartRepo.AddItem(
			ctx,
			cartID,
			productID,
			quantity,
		)
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	log.Printf(
		"user=%d added product=%d quantity=%d",
		userID,
		productID,
		quantity,
	)
	return s.cartRepo.AddItem(
		ctx,
		cart.ID,
		productID,
		quantity,
	)

}

func (s *CartService) GetCart(
	ctx context.Context,
	userID int64,
) ([]model.CartItemInfo, error) {
	cart, err := s.cartRepo.GetByUserID(
		ctx,
		userID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return []model.CartItemInfo{}, nil
	}
	if err != nil {
		return nil, err
	}
	return s.cartRepo.GetCart(
		ctx,
		cart.ID,
	)
}

func (s *CartService) DeleteItem(
	ctx context.Context,
	userID int64,
	itemID int64,
) error {

	cart, err := s.cartRepo.GetByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	log.Printf(
		"user=%d removed item=%d",
		userID,
		itemID,
	)

	return s.cartRepo.DeleteItem(
		ctx,
		itemID,
		cart.ID,
	)
}
