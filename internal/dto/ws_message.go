package dto

import "encoding/json"

type WSMessage struct {
	Type            string `json:"type"`
	SenderID        uint   `json:"senderId"`
	ReceiverID      uint   `json:"receiverId"`
	Body            string `json:"body"`
	AESKeySender    string `json:"aesKeySender,omitempty"`
	AESKeyReceiver  string `json:"aesKeyReceiver,omitempty"`
	MessageID       uint   `json:"messageId,omitempty"`
	State           string `json:"state,omitempty"`
	ExpiredAt       string `json:"expiredAt,omitempty"`
	FileAttachments string `json:"fileAttachments,omitempty"`
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
