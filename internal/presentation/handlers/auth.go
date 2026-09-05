package handlers

import (
	"encoding/json"
	"net/http"
	app_errors "project/internal/application/errors"
	app_context "project/internal/presentation/context"
	"project/internal/presentation/mappers"
	pres_ports "project/internal/presentation/ports"
	"project/internal/presentation/response"
	"project/internal/presentation/schemas"
)

type AuthHandler struct {
	authService pres_ports.AuthService
}

func NewAuthHandler(authService pres_ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary  Регистрация пользователя
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    request body     schemas.CreateUserSchema true "Данные пользователя"
// @Success  201     {object} schemas.UserCreatedSchema
// @Failure  400     {object} schemas.ErrorSchema
// @Router   /auth/register [post]
func (handler *AuthHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	logger := app_context.LoggerFromContext(req.Context())

	defer req.Body.Close()
	var request schemas.CreateUser
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		errorResponse := schemas.Error{Error: "Invalid request body"}
		response.Error(w, http.StatusBadRequest, errorResponse)
		logger.Error("error create user", "error", err)
		return
	}

	result, err := handler.authService.RegisterUser(req.Context(), mappers.FromCreatedSchemaToDTO(&request))
	if err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error create user", "error", err)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.FromCreatedDTOToSchema(result))
}

// @Summary  Авторизация пользователя
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    request body     schemas.AuthorizeSchema true "Учётные данные"
// @Success  200     {object} schemas.TokensSchema
// @Failure  400     {object} schemas.ErrorSchema
// @Failure  401     {object} schemas.ErrorSchema
// @Router   /auth/login [post]
func (handler *AuthHandler) Authorization(w http.ResponseWriter, req *http.Request) {
	logger := app_context.LoggerFromContext(req.Context())

	defer req.Body.Close()
	var request schemas.Authorize
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		errorResponse := schemas.Error{Error: "Invalid request body"}
		response.Error(w, http.StatusBadRequest, errorResponse)
		logger.Error("error authrozation user", "error", err)
		return
	}

	result, err := handler.authService.AuthorizeUser(req.Context(), mappers.FromAuthorizeSchemaToDTO(&request))
	if err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error authrozation user", "error", err)
		return
	}
	responseOptons := []response.Option{
		response.WithRefreshTokenCookie(result.RefreshToken),
	}

	response.JSON(w, http.StatusOK, mappers.FromTokensDTOToSchema(result), responseOptons...)
}

// @Summary     Обновление access-токена
// @Description Выдаёт новый access-токен по refresh-токену из cookie
// @Tags        auth
// @Produce     json
// @Param       refresh_token header   string true "Refresh-токен в cookie" default(refresh_token=<token>)
// @Success     200           {object} schemas.TokensSchema
// @Failure     401           {object} schemas.ErrorSchema
// @Failure     500           {object} schemas.ErrorSchema
// @Router      /auth/refresh [post]
func (handler *AuthHandler) RefreshToken(w http.ResponseWriter, req *http.Request) {
	logger := app_context.LoggerFromContext(req.Context())

	token, err := req.Cookie("refresh_token")
	if err != nil {
		errorResponse := schemas.Error{Error: "Invalid cookie"}
		response.Error(w, http.StatusUnauthorized, errorResponse)
		logger.Error("error while refreshing token", "error", err)
		return
	}

	newToken, err := handler.authService.RefreshToken(req.Context(), token.Value)
	if err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error while refreshing token", "error", err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.FromTokensDTOToSchema(newToken))
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, req *http.Request) {
	logger := app_context.LoggerFromContext(req.Context())

	inputToken, err := req.Cookie("refresh_token")
	if err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(app_errors.ErrInvalidToken)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error while logout", "error", err)
		return
	}

	if err := handler.authService.Logout(req.Context(), inputToken.Value); err != nil {
		status, errorMessage := mappers.FromApplicationToApiError(app_errors.ErrInvalidToken)
		errorResponse := schemas.Error{Error: errorMessage}
		response.Error(w, status, errorResponse)
		logger.Error("error while logout", "error", err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
