package service

import (
	"context"
	"errors"
	"techzone/internal/model"
	"techzone/internal/repository"

	"github.com/jackc/pgx/v5"
)

type CartService struct {
	cartRepo *repository.CartRepository
}

func NewCartService(
	cartRepo *repository.CartRepository,
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

	return s.cartRepo.DeleteItem(
		ctx,
		itemID,
		cart.ID,
	)
}
