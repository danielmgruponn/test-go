package services

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/requests"
	"test-go/internal/response"
)

type messageService struct {
	messageRepo ports.MessageRepository
}

func NewMessageService(messageRepo ports.MessageRepository) ports.MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) SaveMessage(message requests.BodyMessageRequest) (*response.NewMessageResponse, error) {
	messageNew := &domain.Message{
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Body:           message.Content,
		State:          message.Status,
		AESKeySender:   message.AESKeySender,
		AESKeyReceiver: message.AESKeyReceiver,
	}

	newMns, err := s.messageRepo.CreateMessage(messageNew)
	if err != nil {
		return nil, err
	}

	return &response.NewMessageResponse{
		ID:             newMns.ID,
		SenderID:       uint(messageNew.SenderID),
		ReceiverID:     uint(messageNew.ReceiverID),
		Content:        messageNew.Body,
		Status:         messageNew.State,
		ExpiresAt:      "",
		AESKeySender:   messageNew.AESKeySender,
		AESKeyReceiver: messageNew.AESKeyReceiver,
		Event:          "chat",
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
		ID:             newMns.ID,
		SenderID:       uint(newMns.SenderID),
		ReceiverID:     uint(newMns.ReceiverID),
		Content:        newMns.Body,
		Status:         state,
		ExpiresAt:      "",
		AESKeySender:   newMns.AESKeySender,
		AESKeyReceiver: newMns.AESKeyReceiver,
		Event:          "receiver_mns",
	}, nil
}
