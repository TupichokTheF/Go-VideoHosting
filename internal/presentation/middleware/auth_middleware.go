package app_middleware

import (
	"net/http"
	app_ports "project/internal/application/ports"
	"project/internal/presentation/context"
	"project/internal/presentation/response"
	"project/internal/presentation/schemas"
	"strings"
)

func AuthMiddleware(authService app_ports.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			accessToken := req.Header.Get("Authorization")
			if !strings.HasPrefix(accessToken, "Bearer ") {
				errorResponse := schemas.Error{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
			if ok := authService.IsLoggedOut(req.Context(), accessToken); ok {
				errorResponse := schemas.Error{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			userID, ok := authService.IsAuthorized(req.Context(), accessToken)
			if !ok {
				errorResponse := schemas.Error{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			ctx := app_context.ContextWithUserID(req.Context(), userID)
			req = req.WithContext(ctx)

			next.ServeHTTP(w, req)
		})
	}
}
