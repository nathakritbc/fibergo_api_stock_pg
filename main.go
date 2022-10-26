package main

import (
	"fibergo_api_stock_pg/configs"
	"fibergo_api_stock_pg/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(configs.AppConfig)
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("hello ,hello word")
	})

	port := configs.Config("PORT")

	//config for customization
	app.Use(cors.New(configs.CorsConfigDefault))
	app.Use(logger.New())

	configs.InitDatabase()
	routes.SetupRoutes(app)

	app.Listen(":" + port)
}
