package ports

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByNickname(nickname string) (*dto.UserDTO, error)
}
