package routers

import (
	http_handlers "project/internal/presentation/handlers/http"

	"github.com/go-chi/chi/v5"
)

func WithAuthRouter(handler *http_handlers.AuthHandler) Option {
	return func(router chi.Router) {
		router.Route(prefix+"/auth", func(r chi.Router) {
			r.Post("/register", handler.CreateUser)
			r.Post("/login", handler.Authorization)
			r.Post("/refresh", handler.RefreshToken)
			r.Post("/logout", handler.Logout)
		})
	}
}
