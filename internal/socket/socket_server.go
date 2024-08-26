package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"test-go/internal/core/domain"
	"test-go/internal/dto"
	"time"

	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

type SocketHandler struct {
	clients sync.Map
	db      *gorm.DB
}

func NewSocketHandler(db *gorm.DB) *SocketHandler {
	return &SocketHandler{
		clients: sync.Map{},
		db:      db,
	}
}

func (h *SocketHandler) HandleSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {

		userId := c.Locals("id").(string)
		fmt.Printf("User %s connected\n", userId)

		var user domain.User
		if err := h.db.First(&user, userId).Error; err != nil {
			log.Printf("Error finding user: %v\n", err)
			return
		}
		fmt.Printf("User: %v\n", user)

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
				h.handleNewMessage(user.ID, &msg)
			case "status_update":
				h.handleStatusUpdate(user.ID, &msg)
			case "read_receipt":
				h.handleReadReceipt(user.ID, &msg)
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

	fmt.Printf("Expires at: %v\n", t)

	message := domain.Message{
		SenderID:       senderID,
		ReceiverID:     msg.ReceiverID,
		Body:           msg.Body,
		AESKeySender:   msg.AESKeySender,
		AESKeyReceiver: msg.AESKeyReceiver,
		State:          "sent",
		ExpiredAt:      t,
	}
	tx := h.db.Begin()

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating message: %v\n", err)
		return
	}

	var fileUploads []dto.FileUpload
	if err := json.Unmarshal([]byte(msg.FileAttachments), &fileUploads); err != nil {
		tx.Rollback()
		log.Printf("Error unmarshalling file attachments: %v\n", err)
		return
	}

	for _, file := range fileUploads {
		fileAttachment := domain.FileAttachment{
			MessageID: message.ID,
			FileName:  file.FileName,
			FileSize:  file.FileSize,
			FileType:  file.FileType,
			FileURL:   file.FileURL,
		}

		if err := tx.Create(&fileAttachment).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating file attachment: %v\n", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return
	}

	// send message to receiver if online
	if conn, ok := h.clients.Load(msg.ReceiverID); ok {
		wsConn := conn.(*websocket.Conn)
		err := wsConn.WriteJSON(dto.WSMessage{
			Type:            "new_message",
			SenderID:        senderID,
			ReceiverID:      msg.ReceiverID,
			Body:            msg.Body,
			AESKeySender:    msg.AESKeySender,
			AESKeyReceiver:  msg.AESKeyReceiver,
			MessageID:       message.ID,
			Status:          message.State,
			FileAttachments: msg.FileAttachments,
		})

		if err != nil {
			log.Println("Error sending message to receiver:", err)
		} else {
			message.State = "delivered"
			h.db.Save(&message)
		}
	}

	// send confirmation to sender
	confirmMsg := dto.WSMessage{
		Type:      "message_sent",
		MessageID: message.ID,
		Status:    message.State,
	}

	if conn, ok := h.clients.Load(senderID); ok {
		err := conn.(*websocket.Conn).WriteJSON(confirmMsg)
		if err != nil {
			log.Println("Error sending message confirmation to sender:", err)
		}
	}
}

func (h *SocketHandler) handleStatusUpdate(userID uint, msg *dto.WSMessage) {
	var message domain.Message
	if err := h.db.First(&message, msg.MessageID).Error; err != nil {
		log.Println("Error finding message:", err)
		return
	}

	if message.ReceiverID != userID {
		log.Println("User does not have permission to update message status")
		return
	}

	message.State = msg.Status
	if err := h.db.Save(&message).Error; err != nil {
		log.Println("Error updating message status:", err)
		return
	}

	if conn, ok := h.clients.Load(message.SenderID); ok {
		err := conn.(*websocket.Conn).WriteJSON(dto.WSMessage{
			Type:      "status_update",
			MessageID: msg.MessageID,
			Status:    msg.Status,
		})
		if err != nil {
			log.Println("Error sending status update to sender:", err)
		}
	}
}

func (h *SocketHandler) handleReadReceipt(userID uint, msg *dto.WSMessage) {
	var message domain.Message
	if err := h.db.First(&message, msg.MessageID).Error; err != nil {
		log.Println("Error finding message:", err)
		return
	}

	if message.ReceiverID != userID {
		log.Println("User does not have permission to read message")
		return
	}

	message.State = "read"
	if err := h.db.Save(&message).Error; err != nil {
		log.Println("Error updating message status:", err)
		return
	}

	readMsg := dto.WSMessage{
		Type:      "status_update",
		MessageID: msg.MessageID,
		Status:    "read",
	}

	if conn, ok := h.clients.Load(message.SenderID); ok {
		err := conn.(*websocket.Conn).WriteJSON(readMsg)
		if err != nil {
			log.Println("Error sending read receipt to sender:", err)
		}
	}
}
