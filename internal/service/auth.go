package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"techzone/internal/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByLogin(
		ctx context.Context,
		login string,
	) (*model.User, error)

	Create(
		ctx context.Context,
		user *model.User,
	) (int64, error)
}

// инкапсулирует сценарии регистрации и входа
type AuthService struct {
	userRepo UserRepository
}

// создает auth сервис поверх пользовательского репозитория
func NewAuthService(
	userRepo UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// валидирует вход и создает пользователя с захешированным паролем
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

// проверяет логин и пароль и возвращает пользователя для выдачи токена
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
		log.Printf(
			"failed login for %s",
			input.Login,
		)
		return nil, errors.New("invalid login or password")

	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)
	if err != nil {
		log.Printf(
			"failed login for %s",
			input.Login,
		)
		return nil, errors.New("invalid login or password")
	}
	log.Printf(
		"user_id=%d login successful",
		user.ID,
	)
	return user, nil
}
