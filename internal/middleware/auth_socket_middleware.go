package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// Get token from query parameter
			tokenString := c.Query("token")

			// If token is not in query, check Authorization header
			if tokenString == "" {
				authHeader := c.Get("Authorization")
				if authHeader != "" {
					tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

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
			userID := fmt.Sprintf("%.0f", claims["id"].(float64))
			c.Locals("id", userID)
			c.Locals("username", claims["username"])

			fmt.Printf("Authenticated user ID: %s\n", userID)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
