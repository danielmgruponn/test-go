package routes

import (
	"test-go/internal/handlers"
	"test-go/internal/middleware"
	"test-go/internal/socket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupSocketRoutes(app *fiber.App, socketHandler *socket.SocketHandler, webRTCController *handlers.WebRTCHandler) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/auth", middleware.WebSocketAuthMiddleware(), socketHandler.HandleSocket())

	app.Get("/ws/webrtc", middleware.WebSocketAuthMiddleware(), webRTCController.HandlerWebRTC)
}
