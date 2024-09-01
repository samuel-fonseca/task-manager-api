package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samuel-fonseca/task-manager-api/database"
	"github.com/samuel-fonseca/task-manager-api/model"
)

func DeserializeUser(c *fiber.Ctx) error {
	var token string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		token = c.Cookies("token")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No authorization token.",
		})
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claim.",
		})
	}

	var user model.User
	database.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	c.Locals("user", model.FilterUserRecord(&user))

	return c.Next()
}
