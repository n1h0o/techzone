package service

import (
	"context"
	"techzone/internal/model"
	"testing"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	User      *model.User
	GetErr    error
	CreateErr error
}

func (m *MockUserRepo) GetByLogin(
	ctx context.Context,
	login string,
) (*model.User, error) {
	return m.User, m.GetErr
}

func (m *MockUserRepo) Create(
	ctx context.Context,
	user *model.User,
) (int64, error) {
	return 1, m.CreateErr
}

func TestLoginSuccess(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)

	repo := &MockUserRepo{
		User: &model.User{
			ID:           1,
			Login:        "admin",
			PasswordHash: string(hash),
		},
	}

	service := NewAuthService(repo)

	user, err := service.Login(
		context.Background(),
		LoginInput{
			Login:    "admin",
			Password: "123456",
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf(
			"expected user id 1, got %d",
			user.ID,
		)
	}
}

func TestRegister_Success(t *testing.T) {
	repo := &MockUserRepo{
		GetErr: pgx.ErrNoRows,
	}

	service := NewAuthService(repo)

	err := service.Register(
		context.Background(),
		RegisterInput{
			Login:    "newuser",
			Email:    "newuser@mail.ru",
			Password: "123456",
		},
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestRegisterLoginAlreadyExists(t *testing.T) {
	repo := &MockUserRepo{
		User: &model.User{
			ID:    1,
			Login: "admin",
		},
	}

	service := NewAuthService(repo)

	err := service.Register(
		context.Background(),
		RegisterInput{
			Login:    "admin",
			Email:    "admin@mail.ru",
			Password: "123456",
		},
	)
	if err == nil {
		t.Fatal("expected error")
	}
}
