package handlers

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

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
		log.Printf("New WebRTC connection\n")

		log.Printf("Locals: %v\n", ws.Locals("id"))
		userId := ws.Locals("id").(string)
		log.Printf("User %s connected to WebRTC\n", userId)

		id, err := strconv.Atoi(userId)
		if err != nil {
			log.Println("Error converting userId to int:", err)
			return
		}
		userIdInt := uint(id)
		log.Printf("User %s:%d connected to WebRTC\n", userId, userIdInt)
		h.clients.Store(userIdInt, ws)

		defer func() {
			log.Printf("User %s disconnected from WebRTC\n", userId)
			h.clients.Delete(userId)
			ws.Close()
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Printf("read error: %v\n", err)
				break
			}
			log.Printf("Message: %v\n", string(msg))
			var message map[string]interface{}
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("JSON Unmarshal error: %v\n", err)
				continue
			}
			h.handleSignal(userId, message)
		}
	})
}

func (h *WebRTCHandler) handleSignal(userId string, message map[string]interface{}) {
	log.Printf("Handling signal from user %s\n", userId)
	log.Printf("Message: %v\n", message)
	messageType, ok := message["type"].(string)
	if !ok {
		log.Println("Invalid message type")
		return
	}

	to, ok := message["to"].(string)
	if !ok {
		log.Println("Invalid 'to' field")
		return
	}

	// show all clients
	log.Printf("userId: %s", userId)
	log.Printf("to: %s", to)
	log.Println("All clients:")
	h.clients.Range(func(key, value any) bool {
		log.Printf("Client %s", key)
		return true
	})
	// convert to uint
	toUint, err := strconv.Atoi(to)
	if err != nil {
		log.Println("Error converting to uint:", err)
		return
	}
	toUintInt := uint(toUint)
	targetConn, ok := h.clients.Load(toUintInt)
	if !ok {
		log.Printf("Target client %s not found\n", to)
		return
	}

	ws, ok := targetConn.(*websocket.Conn)
	if !ok {
		log.Printf("Invalid connection type for client %s\n", to)
		return
	}

	signalData, err := json.Marshal(message["signal"])
	if err != nil {
		log.Println("Error marshalling signal data: ", err)
		return
	}

	response := map[string]interface{}{
		"type":   messageType,
		"from":   userId,
		"signal": string(signalData),
	}

	err = ws.WriteJSON(response)
	if err != nil {
		log.Printf("Error sending %s message to target: %v\n", messageType, err)
	}
}
