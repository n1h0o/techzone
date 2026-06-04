package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateToken(
	userID int64,
	login string,
	role string,
	secret string,
) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Login:  login,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(
				now,
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	tokenString, err := token.SignedString(
		[]byte(secret),
	)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
