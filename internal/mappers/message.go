package mappers

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

func MapMessageDomainToDTO(message domain.Message) *dto.MessageDTO {
	return &dto.MessageDTO{
		ID:                message.ID,
		SenderID:          message.SenderID,
		ReceiverID:        message.ReceiverID,
		Body:              message.Body,
		State:             message.State,
		AESKeySender:      message.AESKeySender,
		AESKeyReceiver:    message.AESKeyReceiver,
		CreatedAt:         message.CreatedAt,
		ExpiredAt:         message.ExpiredAt,
		NumberAttachments: message.NumberAttachments,
		FileAttachments:   mapFileAttachments(message.FileAttachments),
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

func MapMessagesDomainToDTO(messages []domain.Message) []dto.MessageDTO {
	var dtoMessages []dto.MessageDTO
	for _, message := range messages {
		dtoMessage := dto.MessageDTO{
			ID:                message.ID,
			SenderID:          message.SenderID,
			ReceiverID:        message.ReceiverID,
			Body:              message.Body,
			State:             message.State,
			AESKeySender:      message.AESKeySender,
			AESKeyReceiver:    message.AESKeyReceiver,
			CreatedAt:         message.CreatedAt,
			ExpiredAt:         message.ExpiredAt,
			NumberAttachments: message.NumberAttachments,
			FileAttachments:   mapFileAttachments(message.FileAttachments),
		}
		dtoMessages = append(dtoMessages, dtoMessage)
	}
	return dtoMessages
}
