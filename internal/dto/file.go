package dto

type FileUpload struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
	FileURL  string `json:"fileUrl"`
}

type FileAttachment struct {
	MessageID string `json:"messageId"`
	FileName  string `json:"fileName"`
	FileType  string `json:"fileType"`
	FileSize  int64  `json:"fileSize"`
	FileURL   string `json:"fileUrl"`
}

type NewFileAttachment struct {
	ID string `json:"id"`
}
