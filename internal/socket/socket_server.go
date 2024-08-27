package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

		userId := c.Locals("id").(string)
		// Convert the userId to uint
		id, err := strconv.Atoi(userId)
		if err != nil {
			log.Println("Error converting userId to int:", err)
			return
		}
		userIdInt := uint(id)

		fmt.Printf("User %d connected\n", userIdInt)

		h.clients.Store(userId, c)

		defer func() {
			h.clients.Delete(userId)
			c.Close()
		}()

		for {
			var msg dto.WSMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				break
			}

			fmt.Println("Received message:", msg)

			switch msg.Type {
			case "message":
				h.handleNewMessage(userIdInt, &msg)
			case "status_update":
				h.handleStatusUpdate(userIdInt, &msg)
			case "read_receipt":
				h.handleReadReceipt(userIdInt, &msg)
			default:
				log.Println("Unknown message type:", msg.Type)
			}
		}
	})
}

func (h *SocketHandler) handleNewMessage(senderID uint, msg *dto.WSMessage) {

	t, err := time.Parse(time.RFC3339, msg.ExpiresAt)
	if err != nil {
		return
	}

	var fileUploads []dto.FileUpload
	if err := json.Unmarshal([]byte(msg.FileAttachments), &fileUploads); err != nil {
		log.Printf("Error unmarshalling file attachments: %v\n", err)
		return
	}

	log.Printf("File attachments: %v\n", fileUploads)

	newMessage, err := h.messageHandler.CreateMessage(dto.Message{
		Event:             "message",
		SenderID:          senderID,
		ReceiverID:        msg.ReceiverID,
		Body:              msg.Body,
		AESKeySender:      msg.AESKeySender,
		AESKeyReceiver:    msg.AESKeyReceiver,
		Type:              msg.Type,
		State:             msg.State,
		ExpiresAt:         t,
		NumberAttachments: uint(len(fileUploads)),
	})

	if err != nil {
		log.Println("Error creating message:", err)
		return
	}

	for _, file := range fileUploads {
		fileAttachment := dto.FileAttachment{
			MessageID: newMessage.ID,
			FileName:  file.FileName,
			FileType:  file.FileType,
			FileSize:  file.FileSize,
			FileURL:   file.FileURL,
		}

		newFile, err := h.fileHandler.SaveFile(fileAttachment)
		if err != nil {
			log.Println("Error saving file:", err)
			return
		}
		log.Printf("File saved: %v\n", newFile)
	}

	if conn, ok := h.clients.Load(msg.ReceiverID); ok {
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type:            "new_message",
			SenderID:        senderID,
			ReceiverID:      msg.ReceiverID,
			Body:            msg.Body,
			AESKeySender:    msg.AESKeySender,
			AESKeyReceiver:  msg.AESKeyReceiver,
			MessageID:       newMessage.ID,
			State:           msg.State,
			FileAttachments: msg.FileAttachments,
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}

	// Send confirmation to sender
	if conn, ok := h.clients.Load(senderID); ok {
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type:            "message_sent",
			SenderID:        senderID,
			ReceiverID:      msg.ReceiverID,
			Body:            msg.Body,
			AESKeySender:    msg.AESKeySender,
			AESKeyReceiver:  msg.AESKeyReceiver,
			MessageID:       newMessage.ID,
			State:           "delivered",
			FileAttachments: msg.FileAttachments,
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (h *SocketHandler) handleStatusUpdate(userID uint, msg *dto.WSMessage) {
	err := h.messageHandler.UpdateMessageState(userID, msg.MessageID, msg.State)
	if err != nil {
		log.Println("Error updating message state:", err)
		return
	}

	if conn, ok := h.clients.Load(msg.ReceiverID); ok {
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type:       "status_update",
			SenderID:   userID,
			ReceiverID: msg.ReceiverID,
			MessageID:  msg.MessageID,
			State:      msg.State,
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (h *SocketHandler) handleReadReceipt(userID uint, msg *dto.WSMessage) {
	err := h.messageHandler.UpdateMessageState(userID, msg.MessageID, msg.State)
	if err != nil {
		log.Println("Error updating message state:", err)
		return
	}

	if conn, ok := h.clients.Load(msg.ReceiverID); ok {
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteJSON(dto.WSMessage{
			Type:       "status_update",
			SenderID:   userID,
			ReceiverID: msg.ReceiverID,
			MessageID:  msg.MessageID,
			State:      msg.State,
		}); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
