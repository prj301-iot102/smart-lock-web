package server

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-fuego/fuego"
	_ "github.com/joho/godotenv/autoload"

	"github.com/prj301-iot102/smart-lock-web/backend/internal/config"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/handlers"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/services"
	// middleware "github.com/prj301-iot102/smart-lock-web/backend/internal/server/middlewares"
)

func NewServer() (*fuego.Server, func()) {
	db := database.NewPool()
	jwt := services.NewJwtAuth()

	serverCfg, _ := env.ParseAs[config.Server]()

	server := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf(":%d", serverCfg.Port)),
		// fuego.WithGlobalMiddlewares(middleware.Cors),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				DisableDefaultServer: true,
				DisableMessages:      true,
			}),
		),
	)

	handlers.AuthRoutes(server, db, jwt)
	handlers.UsersRoutes(server, db)
	// handlers.UsersRoutes(server, db)
	handlers.DeviceRoute(server, db)
	handlers.NfcRoute(server, db)
	handlers.AccessLogRoutes(server, db)
	cleanup := func() { db.Close() }

	return server, cleanup
}
