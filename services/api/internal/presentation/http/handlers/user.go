package http_handlers

import (
	"net/http"
	"project/internal/presentation/http/context"
	"project/internal/presentation/http/mappers"
	pres_ports "project/internal/presentation/http/ports"
	"project/internal/presentation/http/response"
	"project/internal/presentation/http/schemas"
)

type UserHandler struct {
	userService pres_ports.UserService
}

func NewUserHandler(userService pres_ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) UserInfo(w http.ResponseWriter, req *http.Request) {
	logger := app_context.LoggerFromContext(req.Context())

	userID, ok := app_context.UserIDFromContext(req.Context())
	if !ok {
		errorResponse := schemas.Error{Error: "Unauthorized"}
		response.Error(w, http.StatusUnauthorized, errorResponse)
		logger.Error("error while user info: unauthorized")
		return
	}

	result, err := handler.userService.GetUserInfo(req.Context(), userID)
	if err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error while user info", "error", err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.FromUserInfoDTOToSchema(result))
}
