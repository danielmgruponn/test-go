package dto

import "time"

type Message struct {
	Event             string    `json:"event"`
	SenderID          uint      `json:"sender_id"`
	ReceiverID        uint      `json:"receiver_id"`
	Body              string    `json:"body"`
	AESKeySender      string    `json:"aes_key_sender,omitempty"`
	AESKeyReceiver    string    `json:"aes_key_receiver,omitempty"`
	Type              string    `json:"type"`
	State             string    `json:"state,omitempty"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	NumberAttachments uint      `json:"number_attachments,omitempty"`
}

type NewMessage struct {
	MessageId int `json:"mns_id"`
}

type UpdateStatusMessage struct {
	Event     string `json:"event"`
	MessageId int    `json:"mns_id"`
}
