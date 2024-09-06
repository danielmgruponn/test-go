package handlers

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Room struct {
	ID      string
	Clients map[string]*websocket.Conn
	mu      sync.Mutex
}

type GroupCallHandler struct {
	Rooms map[string]*Room
	mu    sync.Mutex
}

func NewGroupCallHandler() *GroupCallHandler {
	return &GroupCallHandler{
		Rooms: make(map[string]*Room),
	}
}

func (h *GroupCallHandler) HandlerGroupCall(c *fiber.Ctx) error {
	userId := c.Locals("id").(uint)
	roomId := c.Params("roomId")

	log.Println(userId)
	id := strconv.Itoa(int(userId))

	return websocket.New(func(ws *websocket.Conn) {
		h.mu.Lock()
		if h.Rooms[roomId] == nil {
			h.Rooms[roomId] = &Room{
				ID:      roomId,
				Clients: make(map[string]*websocket.Conn),
			}
		}
		room := h.Rooms[roomId]
		h.mu.Unlock()

		room.mu.Lock()
		room.Clients[id] = ws
		room.mu.Unlock()

		defer func() {
			room.mu.Lock()
			delete(room.Clients, id)
			room.mu.Unlock()
			ws.Close()
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("Message Read Error: ", err)
				break
			}

			var message map[string]interface{}
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Println("Message Unmarshal Error: ", err)
				continue
			}

			h.handleGroupSignal(roomId, id, message)
		}
	})(c)
}

func (h *GroupCallHandler) handleGroupSignal(roomId, userId string, message map[string]interface{}) {
	messageType, ok := message["type"].(string)
	if !ok {
		log.Println("Message Type Invalid")
		return
	}

	room := h.Rooms[roomId]
	if room == nil {
		log.Printf("Room %s not found\n", roomId)
		return
	}

	log.Println(messageType)
	switch messageType {
	case "join":
		h.bradcastToRoom(room, userId, map[string]interface{}{
			"type":   "user-joined",
			"userId": userId,
		})
	case "offer":
		h.bradcastToRoom(room, userId, map[string]interface{}{
			"type":  "offer",
			"from":  userId,
			"offer": message["offer"],
		})
	case "answer":
		to, ok := message["to"].(string)
		if !ok {
			log.Println("Invalid 'to' field answer")
			return
		}
		h.sendToClient(room, to, map[string]interface{}{
			"type":   "answer",
			"from":   userId,
			"answer": message["answer"],
		})
	case "ice-candidate":
		h.bradcastToRoom(room, userId, map[string]interface{}{
			"type":      "ice-candidate",
			"from":      userId,
			"candidate": message["candidate"],
		})
	case "leave":
		h.bradcastToRoom(room, userId, map[string]interface{}{
			"type":   "user-left",
			"userId": userId,
		})
	default:
		log.Printf("Unsupported message type: %s\n", messageType)
	}
}

func (h *GroupCallHandler) bradcastToRoom(room *Room, senderId string, message map[string]interface{}) {
	room.mu.Lock()
	defer room.mu.Unlock()

	for clientId, conn := range room.Clients {
		if clientId != senderId {
			log.Println(clientId)
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to client %s: %v\n", clientId, err)
			}
		}
	}
}

func (h *GroupCallHandler) sendToClient(room *Room, clientId string, message map[string]interface{}) {
	room.mu.Lock()
	defer room.mu.Unlock()

	if conn, ok := room.Clients[clientId]; ok {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error sending message to client %s: %v\n", clientId, err)
		}
	} else {
		log.Printf("Client %s not found in room\n", clientId)
	}
}
