package ports

import (
	"test-go/internal/core/domain"
)

type MessageRepository interface {
	CreateMessage(message *domain.Message) error
	FindById(id string) (*domain.Message, error)
	FindByUserId(id string) ([]domain.Message, error)
	FindBySenderAndReceiverId(senderId, receiverId string) ([]domain.Message, error)
	UpdateStateByMnsId(id string, state string) (*domain.Message, error)
}
