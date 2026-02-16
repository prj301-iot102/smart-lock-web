package database

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	// pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/config"
)

var instance *pgxpool.Pool

func NewPool() *pgxpool.Pool {
	if instance != nil {
		return instance
	}
	ctx := context.Background()

	cfg, _ := env.ParseAs[config.Database]()
	pgxConfig, _ := pgxpool.ParseConfig(cfg.DatabaseURL())
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// pgxUUID.Register(conn.TypeMap())
		return nil
	}

	instance, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
	}

	if err = instance.Ping(ctx); err != nil {
		fmt.Printf("Unable to ping database: %v\n", err)
	}

	return instance
}
