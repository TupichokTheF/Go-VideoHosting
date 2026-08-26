package routers

import (
	"project/internal/presentation/handlers"
	app_middleware "project/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

func WithVideoRouter(handler *handlers.VideoHandler, authService app_middleware.AuthManager) Option {
	return func(router chi.Router) {
		router.Route("/video", func(r chi.Router) {
			r.With(app_middleware.AuthMiddleware(authService)).
				Post("/add", handler.AddVideo)

			r.Get("/get", handler.GetVideo)
		})
	}
}
