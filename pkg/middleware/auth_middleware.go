package middleware

import (
	"fiber/app/models"
	"fiber/config"
	"fiber/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var user struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		tokenString := tokenParts[1]

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		jtiStr, ok := claims["jti"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JTI format"})
		}

		jtiUUID, err := uuid.Parse(jtiStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JTI UUID"})
		}

		var auth models.Auth
		result := config.DB.Where("jti = ?", jtiUUID).First(&auth)
		if result.Error == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is expired"})
		}

		userResult := config.DB.Model(&models.User{}).Select("id,username,email").Where("id = ?", claims["id"]).First(&user)
		if userResult.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		c.Locals("user", user)

		// Proceed to the next handler
		return c.Next()
	}
}
