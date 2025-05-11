package controllers

import (
	"goforum/config"
	"goforum/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c fiber.Ctx) error {
	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Password hashing failed"})
	}

	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

func Login(c fiber.Ctx) error {
	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	config.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Incorrect password"})
	}

	// 生成 JWT token
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token creation failed"})
	}

	return c.JSON(fiber.Map{"token": t})
}

func RequireAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("userRole")
		if role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "Admin access required")
		}
		return c.Next()
	}
}
