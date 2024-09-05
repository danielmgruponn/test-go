package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("AuthMiddleware")
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing auth token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid auth token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		if nickname, ok := claims["nickname"].(string); ok {
			log.Printf("Nickname: %s\n", nickname)
			c.Locals("nickname", nickname)
		}
		if id, ok := claims["id"].(float64); ok {
			log.Printf("ID: %.0f\n", id)
			c.Locals("id", uint(id))
		}

		return c.Next()
	}
}
