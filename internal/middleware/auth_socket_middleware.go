package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("WebSocketAuthMiddleware")
		if websocket.IsWebSocketUpgrade(c) {
			tokenString := c.Query("token")

			if tokenString == "" {
				c.Locals("allowed", false)
				return c.Next()
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil || !token.Valid {
				c.Locals("allowed", false)
				return c.Next()
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Locals("allowed", false)
				return c.Next()
			}

			c.Locals("allowed", true)

			if id, ok := claims["id"].(float64); ok {
				c.Locals("id", uint(id))
			}

			if nickname, ok := claims["nickname"].(string); ok {
				log.Printf("Nickname: %s\n", nickname)
				c.Locals("nickname", nickname)
			}

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
