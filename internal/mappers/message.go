package mappers

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

func MapMessageDomainToDTO(message domain.Message) *dto.Message {
	return &dto.Message{
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Body:           message.Body,
		State:          message.State,
		AESKeySender:   message.AESKeySender,
		AESKeyReceiver: message.AESKeyReceiver,
		ExpiresAt:      message.ExpiredAt,
	}
}

func MapMessageDTOToDomain(message dto.Message) *domain.Message {
	return &domain.Message{
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Body:           message.Body,
		State:          message.State,
		AESKeySender:   message.AESKeySender,
		AESKeyReceiver: message.AESKeyReceiver,
		ExpiredAt:      message.ExpiresAt,
	}
}
