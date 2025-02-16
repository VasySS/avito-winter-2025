package auth

import (
	"fmt"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTGenerator struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewJWTGenerator(secret string, tokenTTL time.Duration) *JWTGenerator {
	return &JWTGenerator{secret: []byte(secret), tokenTTL: tokenTTL}
}

func (g *JWTGenerator) NewAccessToken(user entity.User, currentTime time.Time) (string, error) {
	expirationTime := currentTime.Add(g.tokenTTL)

	token, err := jwt.NewBuilder().
		Expiration(expirationTime).
		Claim("id", user.ID).
		Claim("username", user.Username).
		Build()
	if err != nil {
		return "", fmt.Errorf("token build failed: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, g.secret))
	if err != nil {
		return "", fmt.Errorf("token signing failed: %w", err)
	}

	return string(signed), nil
}
