package routes

import (
	"goforum/controllers"
	"goforum/middleware"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
}

func ProtectedRoutes(app *fiber.App) {
	app.Get("/api/profile", middleware.RequireAuth, func(c fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(fiber.Map{
			"message": "You are authenticated",
			"user":    user,
		})
	})
}
