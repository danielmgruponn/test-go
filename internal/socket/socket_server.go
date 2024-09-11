package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"test-go/internal/dto"
	"test-go/internal/handlers"
	"time"

	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type SocketHandler struct {
	clients        sync.Map
	messageHandler *handlers.MessageHandler
	userHandler    *handlers.UserHandler
	fileHandler    *handlers.FileHandler
}

func NewSocketHandler(messageHandler *handlers.MessageHandler, userHandler *handlers.UserHandler, fileHandler *handlers.FileHandler) *SocketHandler {
	return &SocketHandler{
		clients:        sync.Map{},
		messageHandler: messageHandler,
		userHandler:    userHandler,
		fileHandler:    fileHandler,
	}
}

func (h *SocketHandler) HandleSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		allowed := c.Locals("allowed").(bool)
		log.Printf("Allowed: %v\n", allowed)
		if !allowed {
			log.Println("Not allowed")
			// Send a standardized error message to the frontend indicating the reason for closure
			err := c.WriteJSON(dto.WSMessage{
				Type: dto.WSMessageTypeError,
				Error: &dto.WSError{
					Status:  "401",
					Message: "Invalid token",
				},
			})
			if err != nil {
				log.Printf("Error sending invalid token message: %v\n", err)
			}
			c.Close()
			return
		}

		userId := c.Locals("id").(string)

		h.clients.Store(userId, c)

		defer func() {
			h.clients.Delete(userId)
			c.Close()
		}()

		for {
			var msg dto.WSMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Printf("Read error: %v\n", err)
				break
			}

			log.Println("Received message:", msg)

			switch msg.Type {
			case dto.WSMessageTypeMessage:
				h.handleNewMessage(userId, &msg)
			case dto.WSMessageTypeStatusUpdate:
				h.handleStatusUpdate(userId, &msg)
			case dto.WSMessageTypeReadReceipt:
				h.handleReadReceipt(userId, &msg)
			default:
				log.Println("Unknown message type:", msg.Type)
			}
		}
	})
}

func parseMessageData(data map[string]interface{}) (*dto.MessageData, error) {
	msg := &dto.MessageData{}

	// Safely extract each field, using type assertions
	if val, ok := data["senderId"].(string); ok {
		msg.SenderID = val
	} else {
		return nil, fmt.Errorf("invalid or missing senderId")
	}

	if val, ok := data["receiverId"].(string); ok {
		msg.ReceiverID = val
	} else {
		return nil, fmt.Errorf("invalid or missing receiverId")
	}

	if val, ok := data["body"].(string); ok {
		msg.Body = val
	} else {
		return nil, fmt.Errorf("invalid or missing body")
	}

	if val, ok := data["aesKeySender"].(string); ok {
		msg.AESKeySender = val
	} else {
		return nil, fmt.Errorf("invalid or missing aesKeySender")
	}

	if val, ok := data["aesKeyReceiver"].(string); ok {
		msg.AESKeyReceiver = val
	} else {
		return nil, fmt.Errorf("invalid or missing aesKeyReceiver")
	}

	if val, ok := data["messageId"].(string); ok {
		msg.MessageID = val
	} else {
		return nil, fmt.Errorf("invalid or missing messageId")
	}

	if val, ok := data["state"].(string); ok {
		msg.State = val
	} else {
		return nil, fmt.Errorf("invalid or missing state")
	}

	if val, ok := data["fileAttachments"].(string); ok {
		msg.FileAttachments = val
	} else {
		return nil, fmt.Errorf("invalid or missing fileAttachments")
	}

	if val, ok := data["expiredAt"].(string); ok {
		msg.ExpiredAt = val
	} else {
		return nil, fmt.Errorf("invalid or missing expiredAt")
	}

	return msg, nil
}

func (h *SocketHandler) handleNewMessage(senderID string, msg *dto.WSMessage) {
	log.Println("Handling new message")
	log.Printf("Sender ID: %v\n", senderID)
	log.Printf("Message data: %v\n", msg.Data)

	// data, err := parseMessageData(msg.Data.(map[string]interface{}))
	// if err != nil {
	// 	log.Printf("Error parsing message data: %v\n", err)
	// 	return
	// }

	jsonData, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("Error marshalling message data: %v\n", err)
		return
	}
	var data dto.MessageData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Printf("Error unmarshalling message data: %v\n", err)
		return
	}

	t, err := time.Parse(time.RFC3339, data.ExpiredAt)
	if err != nil {
		log.Printf("Error parsing expired at: %v\n", err)
		return
	}

	var fileUploads []dto.FileUpload
	if err := json.Unmarshal([]byte(data.FileAttachments), &fileUploads); err != nil {
		log.Printf("Error unmarshalling file attachments: %v\n", err)
		return
	}

	log.Printf("File attachments: %v\n", fileUploads)

	err = h.messageHandler.CreateMessage(dto.Message{
		MessageID:         data.MessageID,
		SenderID:          senderID,
		ReceiverID:        data.ReceiverID,
		Body:              data.Body,
		AESKeySender:      data.AESKeySender,
		AESKeyReceiver:    data.AESKeyReceiver,
		State:             data.State,
		ExpiresAt:         t,
		NumberAttachments: uint(len(fileUploads)),
	})

	if err != nil {
		log.Println("Error creating message:", err)
		return
	}

	for _, file := range fileUploads {
		fileAttachment := dto.FileAttachment{
			MessageID: data.MessageID,
			FileName:  file.FileName,
			FileType:  file.FileType,
			FileSize:  file.FileSize,
			FileURL:   file.FileURL,
		}

		_, err := h.fileHandler.SaveFile(fileAttachment)
		if err != nil {
			log.Println("Error saving file:", err)
			return
		}
	}

	if conn, ok := h.clients.Load(data.ReceiverID); ok {
		wsConn := conn.(*websocket.Conn)
		err := wsConn.WriteJSON(dto.WSMessage{
			Type: dto.WSMessageTypeNewMessage,
			Data: dto.MessageData{
				SenderID:        senderID,
				ReceiverID:      data.ReceiverID,
				Body:            data.Body,
				AESKeySender:    data.AESKeySender,
				AESKeyReceiver:  data.AESKeyReceiver,
				MessageID:       data.MessageID,
				State:           data.State,
				FileAttachments: data.FileAttachments,
			},
		})

		if err != nil {
			log.Println("Error sending message to receiver:", err)
			return
		}
	}
}

func (h *SocketHandler) handleStatusUpdate(userID string, msg *dto.WSMessage) {
	log.Printf("Data: %v %T\n", msg.Data, msg.Data)
	// data, ok := msg.Data.(dto.StatusUpdateData)
	// if !ok {
	// 	log.Println("Invalid message data")
	// 	return
	// }

	jsonData, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("Error marshalling message data: %v\n", err)
		return
	}
	var data dto.StatusUpdateData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Printf("Error unmarshalling message data: %v\n", err)
		return
	}

	log.Println("Handling status update")
	log.Printf("Message state: %v\n", data.State)
	log.Printf("Message id: %v\n", data.MessageID)
	err = h.messageHandler.UpdateMessageState(userID, data.MessageID, data.State)
	if err != nil {
		log.Println("Error updating message state:", err)
		return
	}
	log.Printf("Receiver ID: %v %T\n", data.ReceiverID, data.ReceiverID)

	if conn, ok := h.clients.Load(data.ReceiverID); ok {
		log.Println("Sending status update to receiver")
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type: dto.WSMessageTypeStatusUpdate,
			Data: dto.StatusUpdateData{
				SenderID:   userID,
				ReceiverID: data.ReceiverID,
				MessageID:  data.MessageID,
				State:      data.State,
			},
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (h *SocketHandler) handleReadReceipt(userID string, msg *dto.WSMessage) {
	data, ok := msg.Data.(dto.ReadReceiptData)
	if !ok {
		log.Println("Invalid message data")
		return
	}

	log.Println("Handling read receipt")
	log.Printf("Message state: %v\n", data.State)
	log.Printf("Message id: %v\n", data.MessageID)
	err := h.messageHandler.UpdateMessageState(userID, data.MessageID, data.State)
	if err != nil {
		log.Println("Error updating message state:", err)
		return
	}

	if conn, ok := h.clients.Load(data.ReceiverID); ok {
		log.Println("Sending status update to receiver")
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type: dto.WSMessageTypeStatusUpdate,
			Data: dto.ReadReceiptData{
				SenderID:   userID,
				ReceiverID: data.ReceiverID,
				MessageID:  data.MessageID,
				State:      data.State,
			},
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
