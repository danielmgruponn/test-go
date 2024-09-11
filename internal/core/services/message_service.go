package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"test-go/internal/mappers"
)

type messageService struct {
	messageRepo ports.MessageRepository
}

func NewMessageService(messageRepo ports.MessageRepository) ports.MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) SaveMessage(message dto.Message) error {
	messageNew := mappers.MapMessageDTOToDomain(message)

	err := s.messageRepo.CreateMessage(messageNew)
	return err
}

func (s *messageService) GetMyMessages(id string) ([]domain.Message, error) {
	return s.messageRepo.FindByUserId(id)
}

func (s *messageService) UpdateStateMessage(messageId string, state string) error {

	_, err := s.messageRepo.UpdateStateByMnsId(messageId, state)
	if err != nil {
		return err
	}

	return nil
}

func (s *messageService) GetMessagesBySenderAndReceiver(senderID, receiverID string) ([]domain.Message, error) {
	return s.messageRepo.FindBySenderAndReceiverId(senderID, receiverID)
}
