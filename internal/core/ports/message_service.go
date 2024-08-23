package ports

import (
	"test-go/internal/core/domain"
	"test-go/internal/requests"
	"test-go/internal/response"
)

type MessageService interface {
	SaveMessage(message requests.BodyMessageRequest) (*response.NewMessageResponse, error)
	GetMyMessages(id int) ([]domain.Message, error)
	UpdateStateMessage(message requests.UpdateStatusMessage, state string) (*response.NewMessageResponse, error)
}
