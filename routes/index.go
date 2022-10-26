package routes

import (
	"fibergo_api_stock_pg/configs"
	"fibergo_api_stock_pg/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// app.Get("/*", configs.ConfigAuth)
	app.Static("/", "./images")

	// group  route
	apiV1 := app.Group("/api/v1")
	auth := apiV1.Group("/auth") // auth
	product := apiV1.Group("/products", configs.ConfigAuth)
	uploadImage := apiV1.Group("/uploadImages", configs.ConfigAuth)

	user := apiV1.Group("/users", configs.ConfigAuth)

	// auth route
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)

	// user route
	user.Get("/", controllers.DecodeUser)

	// Product route
	product.Get("/", controllers.GetProductsAll)
	product.Get("/:id", controllers.GetProduct)
	product.Post("/", controllers.CreateProduct)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/:id", controllers.DeleteProduct)

	uploadImage.Post("/", controllers.UploadImage)
	uploadImage.Delete("/:fileName", controllers.RemoveImage)

}
