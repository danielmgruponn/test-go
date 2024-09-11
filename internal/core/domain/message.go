package domain

import (
	"time"
)

type Message struct {
	ID                string `gorm:"primary_key"`
	SenderID          string `gorm:"index"`
	ReceiverID        string `gorm:"index"`
	Body              string
	State             string
	AESKeySender      string
	AESKeyReceiver    string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExpiredAt         time.Time
	NumberAttachments uint
	FileAttachments   []FileAttachment `gorm:"foreignKey:MessageID"`
}

func (Message) TableMessages() string {
	return "messages"
}
