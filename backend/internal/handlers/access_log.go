package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/middlewares"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
)

type AccessLogResrouce struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

func (alr *AccessLogResrouce) ListLogs(c fuego.ContextNoBody) ([]database.GetAccessLogsRow, error) {
	ctx := context.Background()
	queries := database.New(alr.db)
	accessLogs, err := queries.GetAccessLogs(ctx)
	if err != nil {
		return []database.GetAccessLogsRow{}, fuego.InternalServerError{
			Detail: "",
		}
	}

	return accessLogs, nil
}

func AccessLogRoutes(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := AccessLogResrouce{
		db:  db,
		jwt: jwt,
	}
	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/accessLog",
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Get(group, "/", rs.ListLogs)
}
