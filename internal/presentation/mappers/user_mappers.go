package mappers

import (
	"project/internal/application/dtos"
	"project/internal/presentation/schemas"
)

func FromUserInfoDTOToSchema(dto *dtos.UserInfo) *schemas.UserInfo {
	return &schemas.UserInfo{
		UserID:    dto.UserID,
		Username:  dto.Username,
		UserEmail: dto.UserEmail,
	}
}
