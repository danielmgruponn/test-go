package domain

type FileAttachment struct {
	ID        uint `gorm:"primary_key"`
	MessageID uint `gorm:"index"`
	FileName  string
	FileType  string
	FileSize  int64
	FileURL   string
}

func (FileAttachment) TableFileAttachment() string {
	return "file_attachments"
}
