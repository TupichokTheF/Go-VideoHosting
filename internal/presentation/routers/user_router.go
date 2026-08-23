package routers

import (
	app_ports "project/internal/application/ports"
	"project/internal/presentation/handlers"
	"project/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

func WithUserRouter(handler *handlers.UserHandler, authService app_ports.AuthService) Option {
	return func(router chi.Router) {
		router.Route("/user", func(r chi.Router) {
			r.Use(app_middleware.AuthMiddleware(authService))

			r.Get("/me", handler.UserInfo)
		})
	}
}
