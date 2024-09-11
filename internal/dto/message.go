package dto

import (
	"time"
)

type Message struct {
	MessageID         string    `json:"messageId"`
	SenderID          string    `json:"senderId"`
	ReceiverID        string    `json:"receiverId"`
	Body              string    `json:"body"`
	AESKeySender      string    `json:"aesKeySender,omitempty"`
	AESKeyReceiver    string    `json:"aesKeyReceiver,omitempty"`
	State             string    `json:"state,omitempty"`
	ExpiresAt         time.Time `json:"expiresAt,omitempty"`
	NumberAttachments uint      `json:"numberAttachments,omitempty"`
}

type MessageDTO struct {
	ID                string           `json:"id"`
	SenderID          string           `json:"senderId"`
	ReceiverID        string           `json:"receiverId"`
	Body              string           `json:"body"`
	State             string           `json:"state"`
	AESKeySender      string           `json:"aesKeySender,omitempty"`
	AESKeyReceiver    string           `json:"aesKeyReceiver,omitempty"`
	CreatedAt         time.Time        `json:"createdAt"`
	ExpiredAt         time.Time        `json:"expiredAt,omitempty"`
	NumberAttachments uint             `json:"numberAttachments,omitempty"`
	FileAttachments   []FileAttachment `json:"fileAttachments,omitempty"`
}

type NewMessage struct {
	MessageId int `json:"mns_id"`
}

type UpdateStatusMessage struct {
	Event     string `json:"event"`
	MessageId int    `json:"mnsId"`
}
