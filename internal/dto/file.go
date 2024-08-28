package dto

type FileUpload struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
	FileURL  string `json:"fileUrl"`
}

type FileAttachment struct {
	MessageID uint   `json:"messageId"`
	FileName  string `json:"fileName"`
	FileType  string `json:"fileType"`
	FileSize  int64  `json:"fileSize"`
	FileURL   string `json:"fileUrl"`
}

type NewFileAttachment struct {
	ID uint `json:"id"`
}
