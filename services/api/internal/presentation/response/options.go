package response

import "net/http"

type Option func(w http.ResponseWriter)

func WithRefreshTokenCookie(token string) Option {
	return func(w http.ResponseWriter) {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
		})
	}
}
