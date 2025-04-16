package controllers

import (
	"goforum/config"
	"goforum/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func CreateComment(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	postId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var commentData struct {
		Content string `json:"content"`
	}
	if err := c.Bind().Body(&commentData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	comment := models.Comment{
		Content: commentData.Content,
		UserID:  user.ID,
		PostID:  uint(postId),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	return c.Status(201).JSON(comment)
}

func GetCommentsByPost(c fiber.Ctx) error {
	postId := c.Params("id")

	var comments []models.Comment
	if err := config.DB.
		Where("post_id = ?", postId).
		Preload("User").
		Order("created_at asc").
		Find(&comments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}

	return c.JSON(comments)
}

func DeleteComment(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	commentId := c.Params("id")

	var comment models.Comment
	if err := config.DB.First(&comment, commentId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Comment not found"})
	}

	if comment.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{"error": "You can only delete your own comments"})
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete comment"})
	}

	return c.JSON(fiber.Map{"message": "Comment deleted successfully"})
}
