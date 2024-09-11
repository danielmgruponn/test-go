package handlers

import (
	"log"
	"sync"
	"test-go/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebRTCHandler struct {
	clients sync.Map
}

func NewWebRTCHandler() *WebRTCHandler {
	return &WebRTCHandler{
		clients: sync.Map{},
	}
}

func (h *WebRTCHandler) HandlerWebRTC() fiber.Handler {
	return websocket.New(func(ws *websocket.Conn) {
		allowed := ws.Locals("allowed").(bool)
		if !allowed {
			log.Printf("Not allowed\n")
			ws.Close()
			return
		}

		userId := ws.Locals("id").(string)

		log.Printf("User %v %T connected to WebRTC\n", userId, userId)
		h.clients.Store(userId, ws)

		defer func() {
			log.Printf("User %v disconnected from WebRTC\n", userId)
			h.clients.Delete(userId)
			ws.Close()
		}()

		for {
			var msg dto.WSRTCMessage
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Printf("Read error: %v\n", err)
				break
			}
			log.Printf("Received message: %s %T\n", msg, msg)
			h.handleSignal(userId, &msg)
		}
	})
}

func (h *WebRTCHandler) handleSignal(userId string, message *dto.WSRTCMessage) {
	log.Printf("Handling signal from %s to %s\n", userId, message.To)
	if conn, ok := h.clients.Load(message.To); ok {
		log.Printf("Sending message to %s\n", message.To)
		wsConn := conn.(*websocket.Conn)
		var response dto.WSRTCMessageResponse
		response.Type = message.Type
		response.From = userId
		response.Signal = message.Signal
		if err := wsConn.WriteJSON(response); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
