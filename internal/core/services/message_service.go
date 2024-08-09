package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
)

type messageService struct {
	messageRepo ports.MessageRepository
}

func NewMessageService(messageRepo ports.MessageRepository) ports.MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) SaveMessage(message *domain.Message) error {
	return s.messageRepo.CreateMessage(message)
}

func (s *messageService) GetMyMessages(userId int) ([]domain.Message, error) {
	return s.messageRepo.FindById(userId)
}
