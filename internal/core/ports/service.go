package ports

import "test-go/internal/dto"

type UserService interface {
	Register(user *dto.RegisterRequest) (uint, error)
	Login(nickname string) (dto.LoginResponse, error)
}
