package service

import (
	"context"
	"errors"
	"strings"
	"techzone/internal/model"
	"techzone/internal/repository"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	input RegisterInput,
) error {

	if input.Login == "" {
		return errors.New("login is required")
	}
	if input.Email == "" || !strings.Contains(input.Email, "@") {
		return errors.New("invalid email")
	}
	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	_, err := s.userRepo.GetByLogin(
		ctx,
		input.Login,
	)
	if err == nil {
		return errors.New("login already exists")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := &model.User{
		Login:        input.Login,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         "client",
	}

	_, err = s.userRepo.Create(
		ctx,
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(
	ctx context.Context,
	input LoginInput,
) (*model.User, error) {
	if input.Login == "" {
		return nil, errors.New("login is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepo.GetByLogin(
		ctx,
		input.Login,
	)
	if err != nil {
		return nil, errors.New("invalid login or password")
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)
	if err != nil {
		return nil, errors.New("invalid login or password")
	}
	return user, nil
}
