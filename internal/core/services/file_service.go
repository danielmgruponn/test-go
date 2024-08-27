package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"time"

	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileService struct {
	s3Client *s3.Client
	fileRepo ports.FileRepository
}

func NewFileService(s3Client *s3.Client, fileRepo ports.FileRepository) *FileService {
	return &FileService{s3Client: s3Client, fileRepo: fileRepo}
}

func (f *FileService) UploadFiles(files []*multipart.FileHeader) ([]dto.FileUpload, error) {
	fileAttachments := make([]dto.FileUpload, 0)
	log.Printf("Uploading %d files", len(files))
	for _, file := range files {
		if file.Size > 10*1024*1024 {
			return nil, fmt.Errorf("file %s is too large", file.Filename)
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

		fileContent, err := file.Open()

		if err != nil {
			return nil, err
		}
		defer fileContent.Close()

		// Upload to S3
		_, err = f.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET")),
			Key:    aws.String(fmt.Sprintf("test/%s", filename)),
			Body:   fileContent,
		})

		if err != nil {
			return nil, err
		}

		client := s3.NewPresignClient(f.s3Client)

		req, err := client.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET")),
			Key:    aws.String(fmt.Sprintf("test/%s", filename)),
		}, func(po *s3.PresignOptions) {
			po.Expires = 24 * time.Hour
		})

		if err != nil {
			log.Printf("Failed to sign request for file %s: %v", file.Filename, err)
			return nil, err
		}

		fileUpload := dto.FileUpload{
			FileName: file.Filename,
			FileSize: file.Size,
			FileType: file.Header.Get("Content-Type"),
			FileURL:  req.URL,
		}

		fileAttachments = append(fileAttachments, fileUpload)
	}
	return fileAttachments, nil
}

func (f *FileService) SaveFile(file *dto.FileAttachment) (dto.NewFileAttachment, error) {
	fileAttachment := domain.FileAttachment{
		MessageID: file.MessageID,
		FileName:  file.FileName,
		FileType:  file.FileType,
		FileSize:  file.FileSize,
		FileURL:   file.FileURL,
	}
	err := f.fileRepo.Create(&fileAttachment)
	if err != nil {
		return dto.NewFileAttachment{}, err
	}
	return dto.NewFileAttachment{
		ID: fileAttachment.ID,
	}, nil
}
