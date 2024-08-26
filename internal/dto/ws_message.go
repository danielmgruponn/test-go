package dto

type WSMessage struct {
	Type            string `json:"type"`
	SenderID        uint   `json:"sender_id"`
	ReceiverID      uint   `json:"receiver_id"`
	Body            string `json:"body"`
	AESKeySender    string `json:"aes_key_sender,omitempty"`
	AESKeyReceiver  string `json:"aes_key_receiver,omitempty"`
	MessageID       uint   `json:"message_id,omitempty"`
	Status          string `json:"status,omitempty"`
	ExpiresAt       string `json:"expires_at,omitempty"`
	FileAttachments string `json:"file_attachments,omitempty"`
}
