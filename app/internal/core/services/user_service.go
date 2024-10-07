package services

import (
	"fmt"
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/pkg/jwt"
)

type userService struct {
    userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *domain.User) error {
    // Aquí deberías hashear la contraseña antes de guardarla
    return s.userRepo.Create(user)
}

func (s *userService) Login(username, password string) (string, error) {
    user, err := s.userRepo.FindByUsername(username)
    if err != nil {
        return "", err
    }
    fmt.Println(password)

    return jwt.GenerateToken(user.ID, user.NickName)
}