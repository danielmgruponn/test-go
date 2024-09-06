package middleware

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

const JWTSecretKeyEnv = "JWT_SECRET_KEY"

type CustomClaims struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

func validateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(JWTSecretKeyEnv)), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func WebSocketAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			tokenString := c.Query("token")

			if tokenString == "" {
				log.Println("Missing token")
			}

			claims, err := validateToken(tokenString)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				c.Locals("allowed", false)
				return websocket.New(func(c *websocket.Conn) {
					c.Close()
				})(c)
			}

			c.Locals("allowed", true)
			c.Locals("id", claims.ID)
			c.Locals("nickname", claims.Nickname)

			log.Printf("Authenticated user: ID=%d, Nickname=%s", claims.ID, claims.Nickname)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
