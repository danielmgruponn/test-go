package ports

import (
	"test-go/internal/dto"
)

type UserService interface {
	Register(user *dto.RegisterRequest) (bool, error)
	Login(nickname string) (dto.LoginResponse, error)
	GetUserById(id string) (dto.UserDTO, error)
	GetUserByNickname(nickname string) (dto.UserDTO, error)
}
