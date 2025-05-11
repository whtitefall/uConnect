package controllers

import (
	"goforum/config"
	"goforum/models"

	"github.com/gofiber/fiber/v3"
)

func CreatePost(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var postData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.Bind().Body(&postData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	post := models.Post{
		Title:   postData.Title,
		Content: postData.Content,
		UserID:  user.ID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.Status(201).JSON(post)
}

func GetPosts(c fiber.Ctx) error {
	var posts []models.Post
	if err := config.DB.Preload("User").Order("created_at desc").Find(&posts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch posts"})
	}
	return c.JSON(posts)
}

func DeletePost(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	postId := c.Params("id")
	userRole := c.Locals("userRole").(string)

	var post models.Post
	if err := config.DB.First(&post, postId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	if post.UserID != user.ID || userRole != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "You can only delete your own posts"})
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete post"})
	}

	return c.JSON(fiber.Map{"message": "Post deleted successfully"})
}
