package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/auth"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResource struct {
	jwt *auth.JwtAuth
	db  *pgxpool.Pool
}

func (ar *AuthResource) Login(c fuego.ContextWithBody[LoginRequest]) (string, error) {
	req, err := c.Body()
	if err != nil {
		return "", fuego.BadRequestError{
			Err:    err,
			Detail: "Invalid login data",
		}
	}

	queries := database.New(ar.db)
	ctx := context.Background()
	user, err := queries.GetAccountByUsername(ctx, req.Username)
	if err != nil {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "Invalid username or password",
		}
	}
	checkPassword, err := auth.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return "", fuego.InternalServerError{
			Detail: "Unable to hash",
		}
	}
	if !checkPassword {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "Invalid password",
		}
	}

	tokenString, err := ar.jwt.MakeJWT(user.ID)
	if err != nil {
		return "", fuego.InternalServerError{
			Err: err,
		}
	}

	return tokenString, nil
}

func AuthRoutes(s *fuego.Server, db *pgxpool.Pool, jwt *auth.JwtAuth) {
	rs := AuthResource{
		db:  db,
		jwt: jwt,
	}

	group := fuego.Group(s, "/auth")
	// fuego.Post(group, "/register", rs.Register)
	fuego.Post(group, "/login", rs.Login)
}
