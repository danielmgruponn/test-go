package main

import (
	"context"
	"log"
	"os"
	"test-go/internal/core/services"
	"test-go/internal/db"
	"test-go/internal/handlers"
	"test-go/internal/repositories"
	"test-go/internal/routes"
	"test-go/internal/socket"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func main() {

	db.Initialize()

	s3Client, err := initS3Client()
	if err != nil {
		log.Fatal("Failed to initialize S3 client")
	}

	// Inicializar repositorios, servicios y handlers
	userRepo := repositories.NewPostgresUserRepository(db.GetDB())
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	mnsRepo := repositories.NewPostgresMessageRepository(db.GetDB())
	mnsService := services.NewMessageService(mnsRepo)
	mnsHandler := handlers.NewMessageHandler(mnsService, userService)

	fileRepo := repositories.NewPostgresFileRepository(db.GetDB())
	fileService := services.NewFileService(s3Client, fileRepo)
	fileHandler := handlers.NewFileHandler(fileService)

	socketHandler := socket.NewSocketHandler(mnsHandler, userHandler, fileHandler)
	webRTC := handlers.NewWebRTCHandler()
	groupWebRTC := handlers.NewGroupCallHandler()

	// Crear la aplicación Fiber
	app := fiber.New()

	routes.SetupSocketRoutes(app, socketHandler, webRTC, groupWebRTC)
	routes.SetupRoutes(app, userHandler, mnsHandler, fileHandler)

	// Iniciar el servidor
	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}

func initS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}
