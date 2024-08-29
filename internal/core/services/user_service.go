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
	userDomain := &domain.User{
		Nickname:   user.Nickname,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}

	error := s.userRepo.Create(userDomain)
	if error != nil {
		return 0, error
	}

	return userDomain.ID, nil
}

func (s *userService) Login(nickname string) (dto.LoginResponse, error) {
	response := dto.LoginResponse{}
	user, err := s.userRepo.FindByNickname(nickname)
	if err != nil {
		return response, err
	}

	jwt, err := jwt.GenerateToken(user.ID, user.Nickname)
	if err != nil {
		return response, err
	}

	response.ID = user.ID
	response.Nickname = user.Nickname
	response.Token = jwt
	response.PrivateKey = user.PrivateKey
	response.PublicKey = user.PublicKey

	return response, nil
}

func (s *userService) GetUserById(id string) (dto.UserDTO, error) {
	response := dto.UserDTO{}
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return response, err
	}

	response.ID = user.ID
	response.Nickname = user.Nickname
	response.PrivateKey = user.PrivateKey
	response.PublicKey = user.PublicKey

	return response, nil
}

func (s *userService) GetUserByNickname(nickname string) (dto.UserDTO, error) {
	response := dto.UserDTO{}
	user, err := s.userRepo.FindByNickname(nickname)
	if err != nil {
		return response, err
	}

	response.ID = user.ID
	response.Nickname = user.Nickname
	response.PrivateKey = user.PrivateKey
	response.PublicKey = user.PublicKey

	return response, nil
}
