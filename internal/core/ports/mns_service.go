package ports

import "test-go/internal/core/domain"

type MessageService interface {
	SaveMessage(message *domain.Message) error
	GetMyMessages(id int) ([]domain.Message, error)
}
