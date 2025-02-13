package jwtlib

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	secret []byte
}

var (
	ErrInvalidToken = errors.New("invalid jwt-token")
)

func NewService(secret string) (*service, error) {
	if secret == "" {

		return nil, errors.New("empty secret")
	}
	return &service{

		secret: []byte(secret),
	}, nil
}

func (s *service) NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to gen new token: %w", err)
	}

	return tokenString, nil
}

func (s *service) ValidateToken(_ context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", errors.Join(ErrInvalidToken, err)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {

		id, ok := claims["uid"].(string)
		if !ok {
			return "", fmt.Errorf("cannot extract user: %w", ErrInvalidToken)
		}

		return id, nil
	}

	return "", ErrInvalidToken
}
