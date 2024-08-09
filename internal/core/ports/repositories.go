package ports

import "test-go/internal/core/domain"

type UserRepository interface {
    Create(user *domain.User) error
    FindByUsername(username string) (*domain.User, error)
}