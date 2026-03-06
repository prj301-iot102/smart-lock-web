package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
)

type AccessLogResrouce struct {
	db *pgxpool.Pool
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

func AccessLogRoutes(s *fuego.Server, db *pgxpool.Pool) {
	rs := AccessLogResrouce{
		db: db,
	}

	group := fuego.Group(s, "/api/accessLog")
	fuego.Get(group, "/", rs.ListLogs)
}
