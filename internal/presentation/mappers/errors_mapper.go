package mappers

import (
	"errors"
	"net/http"
	app_errors "project/internal/application/errors"
	"project/internal/domain/user"
	"project/internal/domain/video"
	"project/internal/infrastructure/storage"
)

func FromApplicationToApiError(appError error) (int, string) {
	var validationErrorUser *user.ValidationError
	var validationErrorVideo *video.ValidationError

	switch {
	case errors.As(appError, &validationErrorUser):
		return http.StatusBadRequest, validationErrorUser.Reason
	case errors.As(appError, &validationErrorVideo):
		return http.StatusBadRequest, validationErrorVideo.Reason
	case errors.Is(appError, user.AlreadyExistError):
		return http.StatusConflict, "User already exists"
	case errors.Is(appError, user.NotFoundError):
		return http.StatusBadRequest, "User wasn't found"
	case errors.Is(appError, user.InvalidPassword):
		return http.StatusUnauthorized, "Invalid password"
	case errors.Is(appError, app_errors.ErrInvalidToken):
		return http.StatusUnauthorized, "Invalid token"
	case errors.Is(appError, app_errors.ErrForbidden):
		return http.StatusForbidden, "Forbidden"
	case errors.Is(appError, video.ErrVideoNotLoaded):
		return http.StatusBadRequest, "Video not loaded"
	case errors.Is(appError, storage.ErrObjectNotFound):
		return http.StatusNotFound, "Object not found"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
