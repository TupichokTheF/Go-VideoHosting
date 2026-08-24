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
				errorResponse := schemas.ErrorSchema{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			refreshToken, err := req.Cookie("refresh_token")
			if err != nil {
				errorResponse := schemas.ErrorSchema{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			if ok := authService.IsLoggedOut(req.Context(), refreshToken.Value); ok {
				errorResponse := schemas.ErrorSchema{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			userID, ok := authService.IsAuthorized(req.Context(), strings.TrimPrefix(accessToken, "Bearer "))
			if !ok {
				errorResponse := schemas.ErrorSchema{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			ctx := app_context.ContextWithUserID(req.Context(), userID)
			req = req.WithContext(ctx)

			next.ServeHTTP(w, req)
		})
	}
}
