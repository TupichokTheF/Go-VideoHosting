package app_middleware

import (
	"log/slog"
	"net/http"
	app_context "project/internal/presentation/context"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		logger := slog.With("request_id", middleware.GetReqID(req.Context()))

		ctx := app_context.ContextWithLogger(req.Context(), logger)
		req = req.WithContext(ctx)

		start := time.Now()
		next.ServeHTTP(w, req)

		logger.Info("request",
			"method", req.Method,
			"path", req.URL.Path,
			"status", ww.Status(),
			"duration", time.Since(start),
		)
	})
}
