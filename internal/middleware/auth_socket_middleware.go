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
		if websocket.IsWebSocketUpgrade(c) {
			// Get token from query parameter
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
				log.Printf("Error parsing token: %v\n", err)
				c.Locals("allowed", false)
				return c.Next()
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Printf("Claims error: %v\n", ok)
				c.Locals("allowed", false)
				return c.Next()
			}

			log.Printf("Claims 123: %v\n", claims)

			c.Locals("allowed", true)
			userID := fmt.Sprintf("%.0f", claims["id"].(float64))
			c.Locals("id", userID)

			log.Printf("Claims 123: %v\n", c.Locals("id"))

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
