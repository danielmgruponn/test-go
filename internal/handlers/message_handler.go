package handlers

import (
	"encoding/json"
	"strconv"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
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

func (h *MessageHandler) CreateMessage(m dto.Message) (*response.NewMessageResponse, error) {

	mns, err := h.messageService.SaveMessage(m)
	if err != nil {
		return nil, err
	}

	fcm, err := services.NewFCMClient()
	if err != nil {
		return nil, err
	}

	data, err := structToStringMap(mns)
	if err != nil {
		return nil, err
	}

	err = fcm.SendMessage("eMwNLXPVQLqGcJnreBgrsE:APA91bFaUdYKZAorS7joDmoapnIPpD4jTKjU_ke5eKYYaIqyO5TB1YGfG6eBaUQKMusdiIM_vdG7ULfBwA6heTTwji4zAKVCJBuyx_W44WGwepUsk2LYHpZjC-KuZy_Mj0coZ9knqZr3", "Test nuevo Mensaje", "Test nuevo Mensaje", data)
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
	messages, err := h.messageService.GetMyMessages(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener mensajes"})
	}
	return c.Status(fiber.StatusOK).JSON(messages)
}

func structToStringMap(inter interface{}) (map[string]string, error) {
	data, err := json.Marshal(inter)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	strMap := make(map[string]string)
	for k, v := range result {
		switch value := v.(type) {
		case string:
			strMap[k] = value
		case float64:
			strMap[k] = strconv.FormatFloat(value, 'f', -1, 64)
		case bool:
			strMap[k] = strconv.FormatBool(value)
		default:
			jsonValue, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			strMap[k] = string(jsonValue)
		}
	}
	return strMap, nil
}

func (h *MessageHandler) UpdateMessageState(userId, messageId uint, messageState string) error {
	return h.messageService.UpdateStateMessage(messageId, messageState)
}
