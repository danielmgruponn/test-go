package main

import (
	"log"
	"os"
	"test-go/internal/core/services"
	"test-go/internal/handlers"
	"test-go/internal/repositories"
	"test-go/internal/socket"
	"test-go/internal/routes"
	"test-go/internal/db"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db.Initialize()
	// Inicializar repositorios, servicios y handlers
	userRepo := repositories.NewPostgresUserRepository(db.GetDB())
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	socketHandler := socket.NewSocketHandler()

	mnsRepo := repositories.NewPostgresMessageRepository(db.GetDB())
	mnsService := services.NewMessageService(mnsRepo)
	mnsHandler := handlers.NewMessageHandler(mnsService)

	webRTC := handlers.NewWebRTCHandler()


	// Crear la aplicación Fiber
	app := fiber.New()

	routes.SetupSocketRoutes(app, socketHandler, webRTC)
	routes.SetupRoutes(app, userHandler, mnsHandler)

	// Iniciar el servidor
	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
