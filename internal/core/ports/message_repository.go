package ports

import "test-go/internal/core/domain"

type MessageRepository interface {
	CreateMessage(message *domain.Message) (*domain.Message, error)
	FindById(id uint) (*domain.Message, error)
	FindByUserId(id uint) ([]domain.Message, error)
	FindBySenderAndReceiverId(senderId, receiverId uint) ([]domain.Message, error)
	UpdateStateByMnsId(id uint, state string) (*domain.Message, error)
}
