package routers

import (
	"log/slog"
	app_middleware "project/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Option func(router chi.Router)
const prefix string = "/api/v1"

func GetRouter(logger *slog.Logger, options ...Option) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(app_middleware.Logging)


	for _, option := range options {
		option(router)
	}

	return router
}
