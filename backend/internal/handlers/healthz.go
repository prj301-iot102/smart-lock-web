package handlers

import (
	"github.com/go-fuego/fuego"
)

type HealthzRespone struct {
	Message string `json:"message"`
}

func HealthzHandler(c fuego.ContextNoBody) (HealthzRespone, error) {
	return HealthzRespone{
		Message: "ok",
	}, nil
}
