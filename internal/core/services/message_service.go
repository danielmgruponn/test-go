package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"test-go/internal/requests"
	"test-go/internal/response"
)

type messageService struct {
	messageRepo ports.MessageRepository
}

func NewMessageService(messageRepo ports.MessageRepository) ports.MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) SaveMessage(message dto.Message) (*response.NewMessageResponse, error) {
	messageNew := &domain.Message{
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Body:           message.Body,
		State:          message.State,
		AESKeySender:   message.AESKeySender,
		AESKeyReceiver: message.AESKeyReceiver,
		ExpiredAt:      message.ExpiresAt,
	}

	newMns, err := s.messageRepo.CreateMessage(messageNew)
	if err != nil {
		return nil, err
	}

	return &response.NewMessageResponse{
		ID: newMns.ID,
	}, nil
}

func (s *messageService) GetMyMessages(userId int) ([]domain.Message, error) {
	return s.messageRepo.FindById(userId)
}

func (s *messageService) UpdateStateMessage(message requests.UpdateStatusMessage, state string) (*response.NewMessageResponse, error) {

	newMns, err := s.messageRepo.UpdateStateByMnsId(message.MessageId, state)
	if err != nil {
		return nil, err
	}

	return &response.NewMessageResponse{
		ID: newMns.ID,
	}, nil
}
