package mappers

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

func MapUserDomainToDTO(user *domain.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:         user.ID,
		NickName:   user.Nickname,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}
}
