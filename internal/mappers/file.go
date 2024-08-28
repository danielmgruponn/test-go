package mappers

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
)

func MapFileDTOToDomain(file dto.FileAttachment) *domain.FileAttachment {
	return &domain.FileAttachment{
		FileName: file.FileName,
		FileType: file.FileType,
		FileSize: file.FileSize,
		FileURL:  file.FileURL,
	}
}

func MapFileDomainToDTO(file domain.FileAttachment) *dto.FileAttachment {
	return &dto.FileAttachment{
		FileName: file.FileName,
		FileType: file.FileType,
		FileSize: file.FileSize,
		FileURL:  file.FileURL,
	}
}

func mapFileAttachments(domainAttachments []domain.FileAttachment) []dto.FileAttachment {
	var dtoAttachments []dto.FileAttachment
	for _, domainAttachment := range domainAttachments {
		dtoAttachment := dto.FileAttachment{
			MessageID: domainAttachment.MessageID,
			FileName:  domainAttachment.FileName,
			FileType:  domainAttachment.FileType,
			FileSize:  domainAttachment.FileSize,
			FileURL:   domainAttachment.FileURL,
		}
		dtoAttachments = append(dtoAttachments, dtoAttachment)
	}
	return dtoAttachments
}
