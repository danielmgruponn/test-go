package socket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"test-go/internal/core/domain"
	"test-go/internal/core/services"
	"test-go/internal/db"
	"test-go/internal/handlers"
	"test-go/internal/repositories"

	// "strconv"
	"sync"
	// "test-go/internal/core/domain"
	// "test-go/internal/core/ports"
	// "test-go/internal/handlers"
	// "test-go/internal/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

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
		// Registrar el nuevo cliente
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
		// Enviar mensaje de bienvenida
		welcomeMsg := MessageObject{
			Event: "welcome",
			Data:  fmt.Sprintf("Bienvenido al chat: %s", userId),
		}
		client.sendMessage(welcomeMsg)
		// Broadcast nuevo usuario conectado
		h.broadcast(MessageObject{
			Event: "newUser",
			Data:  fmt.Sprintf("New user connected: %s", userId),
		}, userId)
		// Loop principal para manejar mensajes
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("error: %v\n", err)
				}
				break
			}
			if messageType == websocket.TextMessage {
				var message MessageObject
				if err := json.Unmarshal(msg, &message); err != nil {
					fmt.Printf("error unmarshalling message: %v\n", err)
					continue
				}
				// Manejar el mensaje según su evento
				switch message.Event {
				case "chat":
					h.handleChatMessage(message)
				// Añade más casos aquí para otros tipos de eventos
				default:
					fmt.Printf("unknown event: %s\n", message.Event)
				}
			}
		}
	})
}

func (h *SocketHandler) handleChatMessage(message MessageObject) {
	if toClient, ok := h.clients[message.To]; ok {
		messageRepo := repositories.NewPostgresMessageRepository(db.GetDB())
		messageService := services.NewMessageService(messageRepo)
		messageHandler := handlers.NewMessageHandler(messageService)

		receiver, err := strconv.Atoi(message.To)
		if err != nil {
			fmt.Println("Error converting string to int: %v", err)
		}
		sender, err := strconv.Atoi(message.From)
		if err != nil {
			fmt.Println("Error converting string to int: %v", err)
		}

		dataMessage := &domain.Message{
			Body:       message.Data,
			ReceiverID: receiver,
			SenderID:   sender,
			State:      "1",
		}
		messageHandler.CreateMessage(dataMessage)
		fmt.Println(dataMessage)
		toClient.sendMessage(message)
	}
}

func (h *SocketHandler) broadcast(message MessageObject, exceptUserId string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userId, client := range h.clients {
		if userId != exceptUserId {
			client.sendMessage(message)
		}
	}
}

func (c *Client) sendMessage(message MessageObject) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if err := c.Conn.WriteJSON(message); err != nil {
		fmt.Printf("error sending message: %v\n", err)
	}
}
