package mappers

import (
	"project/internal/application/dtos"
	"project/internal/presentation/schemas"
)

func FromCreateVideoSchemaToDTO(schema *schemas.CreateVideo) *dtos.CreateVideo {
	return &dtos.CreateVideo{
		OwnerID:     schema.OwnerID,
		Title:       schema.Title,
		Description: schema.Description,
	}
}

func FromGetVideoSchemaToDTO(schema *schemas.GetVideo) *dtos.GetVideo {
	return &dtos.GetVideo{
		VideoID: schema.VideoID,
	}
}

func FromPresignedURLDTOToSchema(dto *dtos.PresignedURL) *schemas.PresignedURL {
	return &schemas.PresignedURL{
		URL: dto.URL,
	}
}
