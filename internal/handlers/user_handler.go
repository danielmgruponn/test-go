package handlers

import (
	"test-go/internal/core/domain"
	"test-go/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
    user := new(domain.User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    if err := h.userService.Register(user); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al registrar usuario"})
    }

    return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
    loginUser := new(domain.User)
    if err := c.BodyParser(loginUser); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }
    token, err := h.userService.Login(loginUser.NickName, loginUser.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
    }

    return c.JSON(fiber.Map{"token": token})
}
