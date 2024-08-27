package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
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
		SenderID:          message.SenderID,
		ReceiverID:        message.ReceiverID,
		Body:              message.Body,
		State:             message.State,
		AESKeySender:      message.AESKeySender,
		AESKeyReceiver:    message.AESKeyReceiver,
		ExpiredAt:         message.ExpiresAt,
		NumberAttachments: message.NumberAttachments,
	}

	newMns, err := s.messageRepo.CreateMessage(messageNew)
	if err != nil {
		return nil, err
	}

	return &response.NewMessageResponse{
		ID: newMns.ID,
	}, nil
}

func (s *messageService) GetMyMessages(id uint) ([]domain.Message, error) {
	return s.messageRepo.FindByUserId(id)
}

func (s *messageService) UpdateStateMessage(messageId uint, state string) error {

	_, err := s.messageRepo.UpdateStateByMnsId(messageId, state)
	if err != nil {
		return err
	}

	return nil
}

func (s *messageService) GetMessagesBySenderAndReceiver(senderID, receiverID uint) ([]domain.Message, error) {
	return s.messageRepo.FindBySenderAndReceiverId(senderID, receiverID)
}
