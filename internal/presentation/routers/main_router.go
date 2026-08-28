package routers

import (
	"github.com/go-chi/chi/v5"
)

type Option func(router chi.Router)
const prefix string = "/api/v1"

func GetRouter(options ...Option) *chi.Mux {
	router := chi.NewRouter()

	for _, option := range options {
		option(router)
	}

	return router
}
