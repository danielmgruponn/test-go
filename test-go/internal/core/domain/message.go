package domain

import "time"

type Message struct {
	ID                uint
	SenderID          int
	ReceiverID        int
	Body              string
	State             string
	AesKeySender      string
	AESKeyReceiver    string
	NumberAttachments int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExpiredAt         time.Time
}

func (Message) TableMessages() string {
	return "messages" // Nombre de la tabla en la base de datos
}
