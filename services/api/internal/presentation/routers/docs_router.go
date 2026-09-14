package routers

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "project/docs"
)

func WithSwagger() Option {
	return func(router chi.Router) {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
	}
}
