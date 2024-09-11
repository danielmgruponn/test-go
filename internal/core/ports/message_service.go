package ports

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

type MessageService interface {
	SaveMessage(message dto.Message) error
	GetMyMessages(id string) ([]domain.Message, error)
	UpdateStateMessage(messageID string, state string) error
	GetMessagesBySenderAndReceiver(senderID, receiverID string) ([]domain.Message, error)
}
