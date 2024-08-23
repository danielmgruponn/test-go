package socket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"test-go/internal/core/services"
	"test-go/internal/db"
	"test-go/internal/handlers"
	"test-go/internal/repositories"
	"test-go/internal/requests"
	"test-go/internal/response"

	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

type SocketHandler struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewSocketHandler() *SocketHandler {
	return &SocketHandler{
		clients: make(map[string]*Client),
	}
}

func (h *SocketHandler) HandleSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userId := c.Locals("id").(string)
		client := &Client{Conn: c}
		h.mu.Lock()
		h.clients[userId] = client
		h.mu.Unlock()
		defer func() {
			h.mu.Lock()
			delete(h.clients, userId)
			h.mu.Unlock()
			c.Close()
		}()
		info := response.NewMessageResponse{
			Event:   "newUser",
			Content: fmt.Sprintf("New user connected: %s", userId),
		}
		client.sendMessage(info)
		h.broadcast(response.NewMessageResponse{
			Event:   "newUser",
			Content: fmt.Sprintf("New user connected: %s", userId),
		}, userId)
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("error: %v\n", err)
				}
				break
			}
			if messageType == websocket.TextMessage {
				var message requests.BodyMessageRequest
				if err := json.Unmarshal(msg, &message); err != nil {
					fmt.Printf("error unmarshalling message: %v\n", err)
					continue
				}
				var stateMns requests.UpdateStatusMessage
				if err := json.Unmarshal(msg, &stateMns); err != nil {
					fmt.Printf("error unmarshalling message: %v\n", err)
					continue
				}
				switch stateMns.Event {
				case "receiver_mns":
					h.handlerStateMnsReceiver(stateMns)
				case "read_mns":
					h.handlerStateMnsRead(stateMns)
				default:
					if message.Event == "chat" {
						h.handleChatMessage(message)
					} else {
						fmt.Printf("unknown event: %s\n", message.Event)
					}
				}
			}
		}
	})
}

func (h *SocketHandler) handleChatMessage(message requests.BodyMessageRequest) {
	if toClient, ok := h.clients[strconv.Itoa(int(message.ReceiverID))]; ok {
		messageRepo := repositories.NewPostgresMessageRepository(db.GetDB())
		messageService := services.NewMessageService(messageRepo)
		messageHandler := handlers.NewMessageHandler(messageService)

		mns, err := messageHandler.CreateMessage(message)
		if err != nil {
			return
		}
		toClient.sendMessage(*mns)
	}
}

func (h *SocketHandler) handlerStateMnsReceiver(message requests.UpdateStatusMessage) {
	messageRepo := repositories.NewPostgresMessageRepository(db.GetDB())
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService)
	mns, err := messageHandler.UpdateStateReceiver(message)
	if err != nil {
		return
	}
	if toClient, ok := h.clients[strconv.Itoa(int(mns.SenderID))]; ok {
		toClient.sendMessage(*mns)
	}
}

func (h *SocketHandler) handlerStateMnsRead(message requests.UpdateStatusMessage) {
	messageRepo := repositories.NewPostgresMessageRepository(db.GetDB())
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService)
	mns, err := messageHandler.UpdateStateRead(message)
	if err != nil {
		return
	}
	if toClient, ok := h.clients[strconv.Itoa(int(mns.SenderID))]; ok {
		toClient.sendMessage(*mns)
	}
}

func (h *SocketHandler) broadcast(message response.NewMessageResponse, exceptUserId string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userId, client := range h.clients {
		if userId != exceptUserId {
			client.sendMessage(message)
		}
	}
}

func (c *Client) sendMessage(message response.NewMessageResponse) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if err := c.Conn.WriteJSON(message); err != nil {
		fmt.Printf("error sending message: %v\n", err)
	}
}
