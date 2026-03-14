package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moznion/go-optional"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/middlewares"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
)

type RoleResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

type SearchRoleRequest struct {
	RoleName optional.Option[string] `json:"role_name"`
}

func (rr *RoleResource) ListRoles(c fuego.ContextNoBody) ([]database.Role, error) {
	ctx := context.Background()
	queries := database.New(rr.db)
	roles, err := queries.ListRoles(ctx)
	if err != nil {
		return []database.Role{}, fuego.InternalServerError{
			Detail: "Unable to get roles",
			Err:    err,
		}
	}

	return roles, nil
}

func (rr *RoleResource) SearchRoleName(c fuego.ContextWithBody[SearchRoleRequest]) ([]database.Role, error) {
	req, err := c.Body()
	if err != nil {
		return []database.Role{}, fuego.BadRequestError{
			Detail: "Invalid body",
			Err:    err,
		}
	}

	ctx := context.Background()
	queries := database.New(rr.db)
	roles, err := queries.SearchRoleName(ctx, req.RoleName)
	if err != nil {
		return []database.Role{}, fuego.InternalServerError{
			Err: err,
		}
	}

	return roles, nil
}

func RoleRoute(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := RoleResource{
		db:  db,
		jwt: jwt,
	}

	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/role",
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()),
	)

	fuego.Get(group, "/", rs.ListRoles)
	fuego.Post(group, "/", rs.SearchRoleName)
}
