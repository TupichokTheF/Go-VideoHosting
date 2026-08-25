package mappers

import (
	"project/internal/application/dtos"
	"project/internal/presentation/schemas"
)

func FromCreatedSchemaToDTO(schema *schemas.CreateUser) *dtos.UserCreate {
	return &dtos.UserCreate{
		UserName:     schema.Username,
		UserPassword: schema.Password,
		UserEmail:    schema.Email,
	}
}

func FromCreatedDTOToSchema(dto *dtos.UserCreated) *schemas.UserCreated {
	return &schemas.UserCreated{
		UserID: dto.UserId,
	}
}

func FromAuthorizeSchemaToDTO(schema *schemas.Authorize) *dtos.Authorize {
	return &dtos.Authorize{
		Username: schema.Username,
		Password: schema.Password,
	}
}

func FromTokensDTOToSchema(dto *dtos.Tokens) *schemas.Token {
	return &schemas.Token{
		AccessToken: dto.AccessToken,
	}
}
