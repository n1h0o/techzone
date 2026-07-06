package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"techzone/internal/app"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {

	server := httptest.NewServer(
		app.NewServer(true).Handler(),
	)

	defer server.Close()

	login := fmt.Sprintf(
		"test_%d",
		time.Now().UnixNano(),
	)

	body := map[string]string{
		"login":    login,
		"email":    login + "@mail.ru",
		"password": "123456",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		server.URL+"/register",
		"application/json", bytes.NewBuffer(data),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Log(err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)

		t.Fatalf(
			"expected %d got %d\nbody=%s",
			http.StatusCreated,
			resp.StatusCode,
			string(body),
		)
	}

	var result map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["message"] != "user created" {
		t.Fatalf(
			"unexpected response: %+v",
			result,
		)
	}
}
