package configs

import "github.com/gofiber/fiber/v2"

var AppConfig = fiber.Config{
	Prefork: false,
}
