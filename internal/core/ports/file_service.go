package ports

import (
	"mime/multipart"
	"test-go/internal/dto"
)

type FileService interface {
	UploadFiles(files []*multipart.FileHeader) ([]dto.FileUpload, error)
	SaveFile(file *dto.FileAttachment) error
}
