package routes

import (
	"goforum/controllers"
	"goforum/middleware"

	"github.com/gofiber/fiber/v3"
)

func CommentRoutes(app *fiber.App) {
	app.Post("/api/posts/:id/comments", middleware.RequireAuth, controllers.CreateComment)
	app.Get("/api/posts/:id/comments", controllers.GetCommentsByPost)
	app.Delete("/api/comments/:id", middleware.RequireAuth, controllers.DeleteComment)
}
