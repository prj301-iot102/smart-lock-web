package server

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-fuego/fuego"
	_ "github.com/joho/godotenv/autoload"

	"github.com/prj301-iot102/smart-lock-web/backend/internal/config"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/handlers"
	// middleware "github.com/prj301-iot102/smart-lock-web/backend/internal/server/middlewares"
)

type Resource struct {
	API handlers.Resource
}

func (rs Resource) NewServer() *fuego.Server {
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

	rs.API.MountRoutes(fuego.Group(server, "/api"))

	return server
}
