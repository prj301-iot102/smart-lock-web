package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/middlewares"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/utils"
)

type UsersResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

func (ur *UsersResource) GetUser(c fuego.ContextNoBody) (database.GetAccountByIdRow, error) {
	user_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return database.GetAccountByIdRow{}, fuego.BadRequestError{
			Detail: "Wrong uuid format",
			Err:    err,
		}
	}

	queries := database.New(ur.db)
	ctx := context.Background()

	user, err := queries.GetAccountById(ctx, user_id)
	if err != nil {
		return database.GetAccountByIdRow{}, fuego.NotFoundError{
			Detail: "User id not exist",
			Err:    err,
		}
	}

	return user, nil
}

func UsersRoutes(s *fuego.Server, db *pgxpool.Pool) {
	rs := UsersResource{
		db:  db,
		jwt: jwt,
	}
	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/users")
	fuego.Use(group, authMiddleware.RequireAuthentication)

	fuego.Get(group, "/{id}", rs.GetUser)
	fuego.Post(group, "/create", rs.CreateUser)
	fuego.Patch(group, "/update", rs.UpdateUserPassword)
}
