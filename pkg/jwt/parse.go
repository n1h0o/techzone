package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(
	tokenString string,
	secret string,
) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
