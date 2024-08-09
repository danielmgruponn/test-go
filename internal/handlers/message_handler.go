package handlers

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MessageHandler struct {
	messageService ports.MessageService
}

func NewMessageHandler(messageService ports.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(m *domain.Message) {
	h.messageService.SaveMessage(m)
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
