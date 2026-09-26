package routers

import (
	http_handlers "project/internal/presentation/handlers/http"
	"project/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

func WithUserRouter(handler *http_handlers.UserHandler, authService app_middleware.AuthManager) Option {
	return func(router chi.Router) {
		router.Route(prefix+"/user", func(r chi.Router) {
			r.Use(app_middleware.AuthMiddleware(authService))

			r.Get("/me", handler.UserInfo)
		})
	}
}
