package handlers

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Room struct {
	ID      string
	Clients sync.Map
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type GroupCallHandler struct {
	Rooms sync.Map
}

func NewGroupCallHandler() *GroupCallHandler {
	return &GroupCallHandler{}
}

func (h *GroupCallHandler) HandleGroupCall(c *fiber.Ctx) error {
	userId := c.Locals("id").(uint)
	roomId := c.Params("roomId")

	id := strconv.FormatUint(uint64(userId), 10)

	return websocket.New(func(ws *websocket.Conn) {
		client := &Client{
			Conn: ws,
			Send: make(chan []byte, 256),
		}

		room := h.getOrCreateRoom(roomId)
		h.addClientToRoom(room, id, client)

		go h.writePump(client)
		h.readPump(room, id, client)
	})(c)
}

func (h *GroupCallHandler) getOrCreateRoom(roomId string) *Room {
	room, _ := h.Rooms.LoadOrStore(roomId, &Room{
		ID: roomId,
	})
	return room.(*Room)
}

func (h *GroupCallHandler) addClientToRoom(room *Room, id string, client *Client) {
	room.Clients.Store(id, client)
}

func (h *GroupCallHandler) removeClientFromRoom(room *Room, id string) {
	if client, ok := room.Clients.LoadAndDelete(id); ok {
		close(client.(*Client).Send)
	}
}

func (h *GroupCallHandler) readPump(room *Room, id string, client *Client) {
	defer func() {
		h.removeClientFromRoom(room, id)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Message Unmarshal Error: ", err)
			continue
		}

		h.handleGroupSignal(room, id, message)
	}
}

func (h *GroupCallHandler) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *GroupCallHandler) handleGroupSignal(room *Room, userId string, message map[string]interface{}) {
	messageType, ok := message["type"].(string)
	if !ok {
		log.Println("Message Type Invalid")
		return
	}

	switch messageType {
	case "join":
		h.broadcastToRoom(room, userId, map[string]interface{}{
			"type":   "user-joined",
			"userId": userId,
		})
	case "offer":
		h.broadcastToRoom(room, userId, map[string]interface{}{
			"type":  "offer",
			"from":  userId,
			"offer": message["offer"],
		})
	case "answer":
		to, ok := message["to"].(string)
		if !ok {
			log.Println("Invalid 'to' field in answer")
			return
		}
		h.sendToClient(room, to, map[string]interface{}{
			"type":   "answer",
			"from":   userId,
			"answer": message["answer"],
		})
	case "ice-candidate":
		h.broadcastToRoom(room, userId, map[string]interface{}{
			"type":      "ice-candidate",
			"from":      userId,
			"candidate": message["candidate"],
		})
	case "leave":
		h.broadcastToRoom(room, userId, map[string]interface{}{
			"type":   "user-left",
			"userId": userId,
		})
	default:
		log.Printf("Unsupported message type: %s\n", messageType)
	}
}

func (h *GroupCallHandler) broadcastToRoom(room *Room, senderId string, message map[string]interface{}) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v\n", err)
		return
	}

	room.Clients.Range(func(key, value interface{}) bool {
		clientId := key.(string)
		client := value.(*Client)
		if clientId != senderId {
			select {
			case client.Send <- jsonMessage:
			default:
				close(client.Send)
				room.Clients.Delete(clientId)
			}
		}
		return true
	})
}

func (h *GroupCallHandler) sendToClient(room *Room, clientId string, message map[string]interface{}) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v\n", err)
		return
	}

	if client, ok := room.Clients.Load(clientId); ok {
		select {
		case client.(*Client).Send <- jsonMessage:
		default:
			close(client.(*Client).Send)
			room.Clients.Delete(clientId)
		}
	} else {
		log.Printf("Client %s not found in room\n", clientId)
	}
}
