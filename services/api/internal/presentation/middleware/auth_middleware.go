package app_middleware

import (
	"context"
	"net/http"
	"project/internal/presentation/context"
	"project/internal/presentation/response"
	"project/internal/presentation/schemas"
	"strings"
)

type AuthManager interface {
	Authenticate(ctx context.Context, accessToken string) (int, bool)
}

func AuthMiddleware(authManager AuthManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			accessToken := req.Header.Get("Authorization")
			if !strings.HasPrefix(accessToken, "Bearer ") {
				errorResponse := schemas.Error{Error: "Unauthorized"}
				response.Error(w, http.StatusUnauthorized, errorResponse)
				return
			}

			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
			userID, ok := authManager.Authenticate(req.Context(), accessToken)
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
