package domain

import "time"

type Message struct {
    ID       	uint   		`json:"id"`
	SenderID    int   		`json:"sender_id"`
	ReceiverID  int   		`json:"receiver_id"`
    Body 		string 		`json:"body"`
    State		string 		`json:"state"`
	AesKey		string 		`json:"aes_key"`
	CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
	ExpiredAt   time.Time   `json:"expired_at"`
}

func (Message) TableMessages() string {
    return "messages" // Nombre de la tabla en la base de datos
}