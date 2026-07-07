package seed

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"

	"techzone/internal/model"
	"techzone/internal/repository"
)

func CreateAdmin(userRepo *repository.UserRepository) error {
	ctx := context.Background()

	_, err := userRepo.GetByLogin(ctx, "admin")
	if err == nil {
		log.Println("admin already exists")
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	_, err = userRepo.Create(ctx, &model.User{
		Login:        "admin",
		Email:        "admin@techzone.local",
		PasswordHash: string(hash),
		Role:         "admin",
	})

	if err != nil {
		return err
	}

	log.Println("admin created")

	return nil
}
