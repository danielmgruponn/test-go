package ports

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
	"test-go/internal/response"
)

type MessageService interface {
	SaveMessage(message dto.Message) (*response.NewMessageResponse, error)
	GetMyMessages(id uint) ([]domain.Message, error)
	UpdateStateMessage(messageID uint, state string) error
}
