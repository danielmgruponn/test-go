package handlers

import (
	"encoding/json"
	"log"
	"strconv"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"test-go/internal/mappers"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	messageService ports.MessageService
	userService    ports.UserService
}

func NewMessageHandler(messageService ports.MessageService, userService ports.UserService) *MessageHandler {
	return &MessageHandler{messageService: messageService, userService: userService}
}

func (h *MessageHandler) CreateMessage(m dto.Message) error {

	err := h.messageService.SaveMessage(m)
	return err

	// fcm, err := services.NewFCMClient()
	// if err != nil {
	// 	return nil, err
	// }

	// data, err := structToStringMap(mns)
	// if err != nil {
	// 	return nil, err
	// }

	// err = fcm.SendMessage("eMwNLXPVQLqGcJnreBgrsE:APA91bFaUdYKZAorS7joDmoapnIPpD4jTKjU_ke5eKYYaIqyO5TB1YGfG6eBaUQKMusdiIM_vdG7ULfBwA6heTTwji4zAKVCJBuyx_W44WGwepUsk2LYHpZjC-KuZy_Mj0coZ9knqZr3", "Test nuevo Mensaje", "Test nuevo Mensaje", data)
	// if err != nil {
	// 	return nil, err
	// }
}

func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	log.Printf("GetMessages")
	id := c.Locals("id").(string)

	partnerID := c.Query("partner_id")
	messages, err := h.messageService.GetMessagesBySenderAndReceiver(id, partnerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener mensajes"})
	}
	messagesDTO := mappers.MapMessagesDomainToDTO(messages)
	return c.Status(fiber.StatusOK).JSON(messagesDTO)
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

func (h *MessageHandler) UpdateMessageState(userId string, messageId string, messageState string) error {
	return h.messageService.UpdateStateMessage(messageId, messageState)
}
