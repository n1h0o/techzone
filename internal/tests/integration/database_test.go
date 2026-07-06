package integration

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func createAdmin(t *testing.T) {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		context.Background(),
		`
		INSERT INTO users(
			login,
			email,
			password_hash,
			role
		)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (login) DO NOTHING
		`,
		"admin",
		"admin@mail.ru",
		string(hash),
		"admin",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func paymentCount(
	t *testing.T,
	orderID int64,
) int {

	t.Helper()

	var count int

	err := db.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM payments
		WHERE order_id = $1
		`,
		orderID,
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	return count
}

func cleanDatabase(t *testing.T) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		`
		TRUNCATE
    		notifications,
   			payments,
    		order_items,
    		orders,
    		cart_items,
    		carts,
    		products,
    		users
		RESTART IDENTITY CASCADE;
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}
