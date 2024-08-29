package mappers

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

func MapUserDomainToDTO(user *domain.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:         user.ID,
		Nickname:   user.Nickname,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}
}

func MapUserDTOToDomain(user *dto.UserDTO) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Nickname:  user.Nickname,
		PublicKey: user.PublicKey,
	}
}

func MapUserDTOToSafeDTO(user *dto.UserDTO) *dto.UserSafeDTO {
	return &dto.UserSafeDTO{
		ID:        user.ID,
		Nickname:  user.Nickname,
		PublicKey: user.PublicKey,
	}
}
