package domain

import "time"

type Message struct {
	ID              uint `gorm:"primary_key"`
	SenderID        uint `gorm:"index"`
	ReceiverID      uint `gorm:"index"`
	Body            string
	State           string
	AESKeySender    string
	AESKeyReceiver  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ExpiredAt       time.Time
	FileAttachments []FileAttachment `gorm:"foreignKey:MessageID"`
}

func (Message) TableMessages() string {
	return "messages" // Nombre de la tabla en la base de datos
}
