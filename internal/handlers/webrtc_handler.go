package handlers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebRTCHandler struct {
	client map[string]*websocket.Conn
	// peerConnection map[string]*webrtc.PeerConnection
	mu sync.Mutex
}

func NewWebRTCHandler() *WebRTCHandler {
	return &WebRTCHandler{
		client: make(map[string]*websocket.Conn),
		// peerConnection: make(map[string]*webrtc.PeerConnection),
	}
}

func (h *WebRTCHandler) HandlerWebRTC(c *fiber.Ctx) error {
    userId := c.Locals("id").(string)

    return websocket.New(func(ws *websocket.Conn) {
        h.mu.Lock()
        h.client[userId] = ws
        h.mu.Unlock()

        defer func() {
            h.mu.Lock()
            delete(h.client, userId)
            h.mu.Unlock()
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

    h.mu.Lock()
    targetConn, ok := h.client[to]
    h.mu.Unlock()

    if !ok {
        fmt.Printf("Target client %s not found\n", to)
        return
    }

    switch messageType {
    case "signal":
        // Reenviar la señal al destinatario
        signalData, err := json.Marshal(message["signal"])
        if err != nil {
            fmt.Println("Error marshalling signal data: ", err)
            return
        }

        err = targetConn.WriteJSON(map[string]interface{}{
            "type":   "signal",
            "from":   userId,
            "signal": string(signalData),
        })
        if err != nil {
            fmt.Println("Error sending signal to target: ", err)
        }
    default:
        fmt.Printf("Unsupported message type: %s\n", messageType)
    }
}
