package ports

import "test-go/internal/dto"

type UserService interface {
	Register(user *dto.RegisterRequest) (uint, error)
	Login(username string) (dto.LoginResponse, error)
}
