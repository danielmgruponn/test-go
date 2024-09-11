package domain

type FileAttachment struct {
	ID        string `gorm:"primary_key"`
	MessageID string `gorm:"index"`
	FileName  string
	FileType  string
	FileSize  int64
	FileURL   string
}

func (FileAttachment) TableFileAttachments() string {
	return "file_attachments"
}
