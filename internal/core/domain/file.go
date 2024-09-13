package domain

import "time"

type FileAttachment struct {
	ID        string `gorm:"primary_key"`
	MessageID string `gorm:"index"`
	FileName  string
	FileType  string
	FileSize  int64
	FileURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (FileAttachment) TableFileAttachments() string {
	return "file_attachments"
}
