package routes

import (
	"test-go/internal/handlers"
	"test-go/internal/middleware"
	"test-go/internal/socket"

	"github.com/gofiber/fiber/v2"
)

func SetupSocketRoutes(app *fiber.App, socketHandler *socket.SocketHandler, webRTCController *handlers.WebRTCHandler) {
	ws := app.Group("/ws")
	ws.Use(middleware.WebSocketAuthMiddleware())
	ws.Get("/", socketHandler.HandleSocket())
	ws.Get("/webrtc", webRTCController.HandlerWebRTC())
}
