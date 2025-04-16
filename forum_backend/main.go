package main

import (
	"goforum/config"
	"goforum/models"
	"goforum/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Post{})

	routes.AuthRoutes(app)
	routes.ProtectedRoutes(app) // 注册需要鉴权的路由
	routes.PostRoutes(app)
	routes.CommentRoutes(app)
	routes.UserRoutes(app)

	// srv := handler.New(generated.NewExecutableSchema(generated.Config{
	// 	Resolvers: &graph.Resolver{},
	// }))
	// 用 adaptor 将 net/http Handler 转为 fiber.Handler
	// 注册 graphql 和 playground
	// app.All("/graphql", adaptor.HTTPHandler(srv))
	// app.All("/playground", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/graphql")))

	app.Listen(":3000")
}
