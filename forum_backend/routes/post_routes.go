package routes

import (
	"goforum/controllers"
	"goforum/middleware"

	"github.com/gofiber/fiber/v3"
)

func PostRoutes(app *fiber.App) {
	app.Post("/api/posts", middleware.RequireAuth, controllers.CreatePost)
	app.Get("/api/posts", controllers.GetPosts)
	app.Delete("/api/posts/:id", middleware.RequireAuth, controllers.DeletePost)
}
