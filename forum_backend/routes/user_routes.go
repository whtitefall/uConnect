package routes

import (
	"goforum/controllers"
	"goforum/middleware"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	app.Get("/api/me", middleware.RequireAuth, controllers.GetMyProfile)
	app.Get("/api/users/:id", controllers.GetUserProfile)
}
