package handlers

import (
	"test-go/internal/core/ports"
	"test-go/internal/requests"
	"test-go/internal/response"
	"test-go/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MessageHandler struct {
	messageService ports.MessageService
}

func NewMessageHandler(messageService ports.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(m requests.BodyMessageRequest) (*response.NewMessageResponse, error) {
	mns, err := h.messageService.SaveMessage(m)
	if err != nil {
		return nil, err
	}
	fcm, err := services.NewFCMClient()
	if err != nil {
		return nil, err
	}
	err = fcm.SendMessage("cipms7TZQoGX-F3enuvDLM:APA91bFkkiwOX_ZvMeOtR4WPEVHPbw10MCqzji164_wbu3sdFlPKyVNb69uapjAzqdvt99sgCX4qNfuCnSkwyR2qvuxj3CWgJWe8gsGFFygFB0OG7SZThy4uFhzDJ-wYOGPsAtrrzm3B", "Test nuevo Mensaje", "Test nuevo Mensaje")
	if err != nil {
		return nil, err
	}
	return mns, nil
}

func (h *MessageHandler) UpdateStateReceiver(m requests.UpdateStatusMessage) (*response.NewMessageResponse, error) {
	mns, err := h.messageService.UpdateStateMessage(m, "2")
	if err != nil {
		return nil, err
	}
	return mns, nil
}

func (h *MessageHandler) UpdateStateRead(m requests.UpdateStatusMessage) (*response.NewMessageResponse, error) {
	mns, err := h.messageService.UpdateStateMessage(m, "3")
	if err != nil {
		return nil, err
	}
	return mns, nil
}

func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	id, ok := user["id"].(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener mensajes"})
	}
	messages, err := h.messageService.GetMyMessages(int(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener mensajes"})
	}
	return c.Status(fiber.StatusOK).JSON(messages)
}
