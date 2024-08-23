package domain

import "time"

type Message struct {
	ID              uint `gorm:"primary_key"`
	SenderID        int
	ReceiverID      int
	Body            string
	State           string
	AesKeySender    string
	AesKeyReceiver  string
	FileAttachments []FileAttachment
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ExpiredAt       time.Time
}

func (Message) TableMessage() string {
	return "messages" // Nombre de la tabla en la base de datos
}
