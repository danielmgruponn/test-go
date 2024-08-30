package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebRTCHandler struct {
	clients sync.Map
}

func NewWebRTCHandler() *WebRTCHandler {
	return &WebRTCHandler{}
}

func (h *WebRTCHandler) HandlerWebRTC(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	log.Printf("User %s connected to WebRTC\n", userId)

	return websocket.New(func(ws *websocket.Conn) {
		h.clients.Store(userId, ws)

		defer func() {
			h.clients.Delete(userId)
			ws.Close()
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("read error: ", err)
				break
			}

			var message map[string]interface{}
			if err := json.Unmarshal(msg, &message); err != nil {
				fmt.Println("JSON Unmarshal error: ", err)
				continue
			}

			h.handleSignal(userId, message)
		}
	})(c)
}

func (h *WebRTCHandler) handleSignal(userId string, message map[string]interface{}) {
	log.Printf("Handling signal from user %s\n", userId)
	log.Printf("Message: %v\n", message)
	messageType, ok := message["type"].(string)
	if !ok {
		fmt.Println("Invalid message type")
		return
	}

	to, ok := message["to"].(string)
	if !ok {
		fmt.Println("Invalid 'to' field")
		return
	}

	targetConn, ok := h.clients.Load(to)
	if !ok {
		fmt.Printf("Target client %s not found\n", to)
		return
	}

	ws, ok := targetConn.(*websocket.Conn)
	if !ok {
		fmt.Printf("Invalid connection type for client %s\n", to)
		return
	}

	signalData, err := json.Marshal(message["signal"])
	if err != nil {
		fmt.Println("Error marshalling signal data: ", err)
		return
	}

	response := map[string]interface{}{
		"type":   messageType,
		"from":   userId,
		"signal": string(signalData),
	}

	err = ws.WriteJSON(response)
	if err != nil {
		fmt.Printf("Error sending %s message to target: %v\n", messageType, err)
	}
}
