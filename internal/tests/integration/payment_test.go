package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"techzone/internal/app"
	"techzone/internal/model"
	"testing"

	"github.com/google/uuid"
)

type loginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type productResponse struct {
	ID int64 `json:"id"`
}

type orderResponse struct {
	OrderID int64 `json:"order_id"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func addToCart(
	t *testing.T,
	url string,
	token string,
	productID int64,
) {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"product_id": productID,
		"quantity":   1,
	})

	req, _ := http.NewRequest(
		http.MethodPost,
		url+"/cart/items",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("add to cart failed: %s", string(data))
	}
}

func createOrder(
	t *testing.T,
	url string,
	token string,
) int64 {

	t.Helper()

	req, _ := http.NewRequest(
		http.MethodPost,
		url+"/orders",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("create order failed: %s", string(data))
	}

	var result orderResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	return result.OrderID
}

func pay(
	t *testing.T,
	url string,
	token string,
	orderID int64,
	key string,
) *model.Payment {

	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"order_id": orderID,
	})

	req, _ := http.NewRequest(
		http.MethodPost,
		url+"/payments",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("payment failed: %s", string(data))
	}

	var payment model.Payment

	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		t.Fatal(err)
	}

	return &payment
}

func register(
	t *testing.T,
	url string,
	login string,
	email string,
	password string,
) {

	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"login":    login,
		"email":    email,
		"password": password,
	})

	resp, err := http.Post(
		url+"/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("register failed: %s", string(data))
	}
}

func login(
	t *testing.T,
	url string,
	login string,
	password string,
) string {

	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})

	resp, err := http.Post(
		url+"/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("login failed: %s", string(data))
	}

	var result loginResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	return result.Token
}

func createProduct(
	t *testing.T,
	url string,
	token string,
) int64 {

	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name":        fmt.Sprintf("product-%s", uuid.NewString()),
		"description": "description",
		"price":       1000,
		"stock":       10,
		"image_url":   "",
	})

	req, _ := http.NewRequest(
		http.MethodPost,
		url+"/products",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("create product failed: %s", string(data))
	}

	var result productResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	return result.ID
}

func TestPaymentIdempotency(t *testing.T) {

	cleanDatabase(t)

	server := httptest.NewServer(
		app.NewServer(true).Handler(),
	)
	defer server.Close()

	//admin

	adminToken := login(
		t,
		server.URL,
		"admin",
		"123456",
	)

	productID := createProduct(
		t,
		server.URL,
		adminToken,
	)

	//user

	loginName := uuid.NewString()

	register(
		t,
		server.URL,
		loginName,
		loginName+"@mail.ru",
		"123456",
	)

	userToken := login(
		t,
		server.URL,
		loginName,
		"123456",
	)

	addToCart(
		t,
		server.URL,
		userToken,
		productID,
	)

	orderID := createOrder(
		t,
		server.URL,
		userToken,
	)

	key := uuid.NewString()

	first := pay(
		t,
		server.URL,
		userToken,
		orderID,
		key,
	)

	second := pay(
		t,
		server.URL,
		userToken,
		orderID,
		key,
	)

	if first.ID != second.ID {
		t.Fatal("payment ids are different")
	}

	if first.TransactionID != second.TransactionID {
		t.Fatal("transaction ids are different")
	}

	count := paymentCount(
		t,
		orderID,
	)

	if count != 1 {
		t.Fatalf(
			"expected 1 payment got %d",
			count,
		)
	}
}
