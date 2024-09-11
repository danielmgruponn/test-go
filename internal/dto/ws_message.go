package dto

import (
	"encoding/json"
)

type WSMessageType string

const (
	WSMessageTypeError        WSMessageType = "error"
	WSMessageTypeNewMessage   WSMessageType = "new_message"
	WSMessageTypeMessageSent  WSMessageType = "message_sent"
	WSMessageTypeStatusUpdate WSMessageType = "status_update"
	WSMessageTypeReadReceipt  WSMessageType = "read_receipt"
	WSMessageTypeMessage      WSMessageType = "message"
)

type WSError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type WSMessage struct {
	Type  WSMessageType `json:"type"`
	Data  interface{}   `json:"data"`
	Error *WSError      `json:"error,omitempty"`
}

type MessageData struct {
	SenderID        string `json:"senderId"`
	ReceiverID      string `json:"receiverId"`
	Body            string `json:"body"`
	AESKeySender    string `json:"aesKeySender"`
	AESKeyReceiver  string `json:"aesKeyReceiver"`
	MessageID       string `json:"messageId"`
	State           string `json:"state"`
	FileAttachments string `json:"fileAttachments"`
	ExpiredAt       string `json:"expiredAt"`
}

type StatusUpdateData struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	MessageID  string `json:"messageId"`
	State      string `json:"state"`
}

type ReadReceiptData struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	MessageID  string `json:"messageId"`
	State      string `json:"state"`
}

type WSRTCMessage struct {
	Type   string          `json:"type"`
	To     string          `json:"to"`
	Signal json.RawMessage `json:"signal"`
}

type WSRTCMessageResponse struct {
	Type   string          `json:"type"`
	From   string          `json:"from"`
	Signal json.RawMessage `json:"signal"`
}
