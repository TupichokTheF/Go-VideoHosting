package routers

import (
	http_handlers "project/internal/presentation/handlers/http"
	app_middleware "project/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

func WithVideoRouter(handler *http_handlers.VideoHandler, authService app_middleware.AuthManager) Option {
	return func(router chi.Router) {
		router.Route(prefix+"/video", func(r chi.Router) {
			r.With(app_middleware.AuthMiddleware(authService)).
				Post("/add", handler.AddVideo)
			r.With(app_middleware.AuthMiddleware(authService)).
				Post("/{video_id}/complete", handler.Complete)
			r.Get("/{video_id}", handler.GetVideo)
		})
	}
}
