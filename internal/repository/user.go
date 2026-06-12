package repository

import (
	"context"
	"techzone/internal/model"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(
	db DBTX,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByLogin(
	ctx context.Context,
	login string,
) (*model.User, error) {

	var user model.User

	err := r.db.QueryRow(
		ctx,
		`SELECT id, login, email, password_hash, role, created_at
		FROM users
		WHERE login = $1
		`,
		login,
	).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		ctx,
		`
			SELECT id, login, email, password_hash, role,created_at
			FROM users
			WHERE id = $1
			`,
		id,
	).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *model.User,
) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO users(
			login,
			email,
			password_hash,
			role
		)
		VALUES ($1,$2,$3,$4)
		RETURNING id
		`,
		user.Login,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&id)

	return id, err
}
