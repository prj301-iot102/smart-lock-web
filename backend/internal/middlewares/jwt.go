package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
)

type key string

const (
	authrizaion           string = "Authorization"
	bearer                string = "Bearer "
	AuthorizationTokenKey key    = "token"
)

type AuthMiddleware struct {
	jwtService *token.JwtAuth
}

func NewAuthMiddleware(jwtService *token.JwtAuth) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authrizaion)
		if authHeader == "" {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Title: "Missing authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)
		token, err := m.jwtService.ValidateJWT(tokenString)
		if err != nil {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Title: "Invalid authorization token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), AuthorizationTokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
