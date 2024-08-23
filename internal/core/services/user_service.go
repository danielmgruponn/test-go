package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"test-go/pkg/jwt"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *dto.RegisterRequest) (uint, error) {
	// Map user request to domain user

	userDomain := &domain.User{
		Nickname:   user.NickName,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}

	error := s.userRepo.Create(userDomain)
	if error != nil {
		return 0, error
	}

	return userDomain.ID, nil
}

func (s *userService) Login(username string) (dto.LoginResponse, error) {
	response := dto.LoginResponse{}
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return response, err
	}

	jwt, err := jwt.GenerateToken(user.ID, user.Nickname)
	if err != nil {
		return response, err
	}

	response.ID = user.ID
	response.NickName = user.Nickname
	response.Token = jwt
	response.PrivateKey = user.PrivateKey
	response.PublicKey = user.PublicKey

	return response, nil
}
