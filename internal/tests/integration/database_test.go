package integration

import (
	"context"
	"testing"
)

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
		TRUNCATE notifications,
		         payments,
		         orders,
		         carts
		RESTART IDENTITY CASCADE;
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}
