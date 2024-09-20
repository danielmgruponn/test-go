package routes

import (
	"test-go/internal/handlers"
	"test-go/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, userController *handlers.UserHandler, mnsController *handlers.MessageHandler, fileController *handlers.FileHandler) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000/",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/", userController.Home)

	app.Post("/register", userController.Register)
	app.Post("/login", userController.Login)

	// Grupo de rutas protegidas
	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Ejemplo de ruta protegida
	api.Get("/users/:id", userController.GetUserById)

	api.Get("/messages", mnsController.GetMessages)

	api.Post("/upload-files", fileController.UploadFiles)

}
