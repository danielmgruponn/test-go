package ports

import "test-go/internal/core/domain"

type MessageRepository interface {
	CreateMessage(message *domain.Message) error
	FindById(id int) ([]domain.Message, error)
}
