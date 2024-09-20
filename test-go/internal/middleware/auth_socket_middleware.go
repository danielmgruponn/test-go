package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
				return websocket.New(func(ws *websocket.Conn) {
					ws.WriteJSON(map[string]interface{}{
							"type": "error",
							"code": 404,
						})
					ws.Close()
				})(c)
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				c.Locals("allowed", false)
				return c.Next()
			}
			c.Locals("allowed", true)
			userID := fmt.Sprintf("%.0f", claims["id"].(float64))
			fmt.Printf("Authenticated user ID: %s\n", userID)
			c.Locals("id", userID)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
