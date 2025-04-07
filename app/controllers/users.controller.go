package controllers

import (
	"fiber/app/models"
	"fiber/config"
	"fiber/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignupUser(c *fiber.Ctx) error {
	var user models.User

	// Parse request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user.Password = utils.GeneratePassword(user.Password)

	result := config.DB.Create(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User signed up successfully"})
}

func LoginUser(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !utils.ComparePassword(input.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	tokenData := utils.AuthDto{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Jti:      uuid.New(),
	}
	token, _ := utils.CreateToken(tokenData)
	refreshToken, _ := utils.CreateRefreshToken(tokenData)
	return c.JSON(fiber.Map{"message": "Login successful", "accessToken": token, "refreshToken": refreshToken})
}

func GetUserDetail(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.JSON(fiber.Map{"status": 200, "message": "User Details Get", "data": user})
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	claims, err := utils.VerifyToken(tokenParts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token structure"})
	}
	jtiString, ok := claims["jti"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid uuid structure"})
	}
	jti, _ := uuid.Parse(jtiString)
	authData := models.Auth{
		Jti: jti,
	}
	result := config.DB.Create(&authData)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failer to logout"})
	}
	return c.JSON(fiber.Map{"status": 200, "message": "User Logout Successful."})
}

func Refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	claims, err := utils.VerifyRefreshToken(tokenParts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	jtiString, ok := claims["jti"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid uuid"})
	}
	jti, _ := uuid.Parse(jtiString)
	authData := models.Auth{
		Jti: jti,
	}
	result := config.DB.Create(&authData)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failer to refresh the token"})
	}
	idStr, _ := claims["id"].(string)
	id, _ := uuid.Parse(idStr)
	email, _ := claims["email"].(string)
	username, _ := claims["username"].(string)

	authDto := utils.AuthDto{
		ID:       id,
		Email:    email,
		Username: username,
		Jti:      uuid.New(),
	}

	token, _ := utils.CreateToken(authDto)
	refreshToken, _ := utils.CreateRefreshToken(authDto)

	return c.JSON(fiber.Map{"status": 200, "message": "User Logout Successful.", "accessToken": token, "refreshToken": refreshToken})
}
