package ports

import "test-go/internal/core/domain"

type FileRepository interface {
	Create(file *domain.FileAttachment) error
}
