package requests

type BodyMessageRequest struct {
	Event           string `json:"event"`
	SenderID        string `json:"sender_id"`
	ReceiverID      string `json:"receiver_id"`
	Body            string `json:"body"`
	AESKeySender    string `json:"aes_key_sender,omitempty"`
	AESKeyReceiver  string `json:"aes_key_receiver,omitempty"`
	Type            string `json:"type"`
	Status          string `json:"status,omitempty"`
	ExpiresAt       string `json:"expires_at,omitempty"`
	FileAttachments int    `json:"file_attachments,omitempty"`
}
