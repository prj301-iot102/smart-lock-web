package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/prj301-iot102/smart-lock-web/backend/internal/config"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "smart-lock-web-access"
)

type JwtAuth struct {
	cfg config.Jwt
}

func NewJwtAuth() *JwtAuth {
	cfg, _ := env.ParseAs[config.Jwt]()

	return &JwtAuth{cfg}
}

func (a *JwtAuth) MakeJWT(
	userID uuid.UUID,
) (string, error) {
	signingKey := a.cfg.JwtSecret
	expiresIn := time.Duration(a.cfg.ExpiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func (a *JwtAuth) ValidateJWT(tokenString string) (uuid.UUID, error) {
	tokenSecret := a.cfg.JwtSecret
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) { return tokenSecret, nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
