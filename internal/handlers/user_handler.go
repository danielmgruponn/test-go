package handlers

import (
	"log"
	"test-go/internal/core/ports"
	"test-go/internal/dto"
	"test-go/internal/mappers"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	created, err := h.userService.Register(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al registrar usuario"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"created": created})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	loginUser := new(dto.LoginRequest)
	if err := c.BodyParser(loginUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	response, err := h.userService.Login(loginUser.Nickname)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	user, err := h.userService.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	u := mappers.MapUserDTOToSafeDTO(&user)

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *UserHandler) GetUserByNickname(c *fiber.Ctx) error {
	nickname := c.Query("nickname")
	log.Printf("Nickname: %v %T", nickname, nickname)

	if nickname == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nickname inválido"})
	}

	user, err := h.userService.GetUserByNickname(nickname)
	if err != nil {
		log.Printf("Error getting the user info by nickname: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	u := mappers.MapUserDTOToSafeDTO(&user)
	log.Printf("User: %v", u)

	return c.Status(fiber.StatusOK).JSON(u)
}
