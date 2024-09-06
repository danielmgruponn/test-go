package handlers

import (
	"log"
	"strconv"
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
		userId := ws.Locals("id").(uint)

		log.Printf("User %d connected to WebRTC\n", userId)
		h.clients.Store(userId, ws)

		defer func() {
			log.Printf("User %d disconnected from WebRTC\n", userId)
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
			log.Printf("Received message Type: %s To: %s\n", msg.Type, msg.To)
			h.handleSignal(userId, &msg)
		}
	})
}

func (h *WebRTCHandler) handleSignal(userId uint, message *dto.WSRTCMessage) {
	id, err := strconv.ParseUint(message.To, 10, 32)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	if conn, ok := h.clients.Load(uint(id)); ok {
		log.Printf("Sending message to %s\n", message.To)
		wsConn := conn.(*websocket.Conn)
		var response dto.WSRTCMessageResponse
		response.Type = message.Type
		response.From = strconv.Itoa(int(userId))
		response.Signal = message.Signal
		if err := wsConn.WriteJSON(response); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
