package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/application/dtos"
	app_context "project/internal/presentation/context"
	"project/internal/presentation/mappers"
	pres_ports "project/internal/presentation/ports"
	"project/internal/presentation/response"
	"project/internal/presentation/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type VideoHandler struct {
	videoService pres_ports.VideoService
}

func NewVideoHandler(videoService pres_ports.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (handler *VideoHandler) AddVideo(w http.ResponseWriter, req *http.Request) {
	userID, ok := app_context.UserIDFromContext(req.Context())
	if !ok {
		errorResponse := schemas.Error{Error: "Unauthorized"}
		response.Error(w, http.StatusUnauthorized, errorResponse)
		return
	}

	var request schemas.CreateVideo = schemas.CreateVideo{OwnerID: userID}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		errorResponse := schemas.Error{Error: "Invalid request body"}
		response.Error(w, http.StatusBadRequest, errorResponse)
		return
	}

	result, err := handler.videoService.CreateVideo(req.Context(), mappers.FromCreateVideoSchemaToDTO(&request))
	if err != nil {
		fmt.Println(err)
		status, messsage := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: messsage}
		response.Error(w, status, errorResponse)
		return
	}

	response.JSON(w, http.StatusOK, mappers.FromPresignedURLDTOToSchema(result))
}

func (handler *VideoHandler) GetVideo(w http.ResponseWriter, req *http.Request) {
	var request schemas.GetVideo
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		errorResponse := schemas.Error{Error: "Invalid request body"}
		response.Error(w, http.StatusBadRequest, errorResponse)
		return
	}

	result, err := handler.videoService.GetVideo(req.Context(), mappers.FromGetVideoSchemaToDTO(&request))
	if err != nil {
		fmt.Println(err)
		status, message := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: message}
		response.Error(w, status, errorResponse)
		return
	}

	response.JSON(w, http.StatusOK, mappers.FromPresignedURLDTOToSchema(result))
}

func (handler *VideoHandler) Complete(w http.ResponseWriter, req *http.Request) {
	userID, ok := app_context.UserIDFromContext(req.Context())
	if !ok {
		errorResponse := schemas.Error{Error: "Unauthorizaed"}
		response.Error(w, http.StatusUnauthorized, errorResponse)
		return
	}

	videoID, err := uuid.Parse(chi.URLParam(req, "video_id"))
	if err != nil {
		errorResponse := schemas.Error{Error: "Bad request"}
		response.Error(w, http.StatusBadRequest, errorResponse)
		return
	}

	err = handler.videoService.CompleteVideo(req.Context(), &dtos.CompleteVideo{VideoID: videoID, UserID: userID})
	if err != nil {
		status, message := mappers.FromApplicationToApiError(err)
		errorResponse := schemas.Error{Error: message}
		response.Error(w, status, errorResponse)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
