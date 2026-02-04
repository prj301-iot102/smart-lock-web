package handlers

import (
	"github.com/go-fuego/fuego"
)

type Resource struct {
}

func (rs Resource) MountRoutes(s *fuego.Server) {
	fuego.Get(s, "/healthz", HealthzHandler)
}
