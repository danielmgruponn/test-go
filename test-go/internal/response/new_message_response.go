package response

type NewMessageResponse struct {
	ID              uint
	SenderID        uint
	ReceiverID      uint
	Content         string
	Status          string
	ExpiresAt       string
	AESKeySender    string
	AESKeyReceiver  string
	Event			string
}