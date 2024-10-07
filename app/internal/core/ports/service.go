package ports

import "test-go/internal/core/domain"

type UserService interface {
    Register(user *domain.User) error
    Login(username, password string) (string, error)
}