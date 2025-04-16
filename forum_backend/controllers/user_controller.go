package controllers

import (
	"goforum/config"
	"goforum/models"

	"github.com/gofiber/fiber/v3"
)

func GetMyProfile(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var fullUser models.User
	if err := config.DB.
		Preload("Posts").
		Preload("Comments").
		First(&fullUser, user.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user profile"})
	}

	return c.JSON(fullUser)
}

func GetUserProfile(c fiber.Ctx) error {
	userId := c.Params("id")

	var user models.User
	if err := config.DB.
		Preload("Posts").
		Preload("Comments").
		First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}
