package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/services"
)

const (
	authrizaion           string = "Authorization"
	bearer                string = "Bearer "
	AuthorizationTokenKey string = "token"
)

func RequireAuthentication(next http.Handler) http.Handler {
	jwtService := services.NewJwtAuth()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authrizaion)
		if authHeader == "" {
			fuego.SendJSONError(w, nil, fuego.BadRequestError{
				Title: "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)

		token, err := jwtService.ValidateJWT(tokenString)
		if err != nil {
			fuego.SendJSONError(w, nil, fuego.BadRequestError{
				Title: "Invalid authorization token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), AuthorizationTokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
