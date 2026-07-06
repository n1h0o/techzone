package handler

import (
	"encoding/json"
	"net/http"
	"techzone/internal/config"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(
	authService *service.AuthService,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Register godoc
//
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterInput true "Данные пользователя"
// @Success 201 {object} handler.MessageResponse
// @Failure 400 {string} string
// @Router /register [post]
func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Only POST method allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}
	var req service.RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authService.Register(
		r.Context(),
		req,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(
		MessageResponse{
			Message: "user created",
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Login godoc
//
// @Summary Авторизация
// @Description Выполняет вход пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginInput true "Логин и пароль"
// @Success 200 {object} handler.LoginResponse
// @Failure 401 {string} string
// @Router /login [post]
func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	var req service.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(
		r.Context(),
		req,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := jwt.GenerateToken(
		user.ID,
		user.Login,
		user.Role,
		h.cfg.JWTSecret,
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(
		LoginResponse{
			Message: "login successful",
			Token:   token,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
