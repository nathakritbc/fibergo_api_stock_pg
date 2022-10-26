package controllers

import (
	"fibergo_api_stock_pg/database"
	"fibergo_api_stock_pg/helpers"
	_ "fibergo_api_stock_pg/helpers"
	"fibergo_api_stock_pg/models"

	"log"

	"github.com/gofiber/fiber/v2"
)

func GetProductsAll(c *fiber.Ctx) error {
	product := []models.Product{}
	db := database.DBConn
	err := db.Find(&product).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product := models.Product{}
	db := database.DBConn
	err := db.Find(&product, id).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	if product.ID == 0 {
		log.Print(product)
		return c.JSON(fiber.Map{})
	}

	return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {

	product := models.Product{}

	err := c.BodyParser(&product)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	errors := helpers.ValidateStructProduct(product)
	if errors != nil {
		log.Print(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errors})
	}

	db := database.DBConn
	err = db.Create(&product).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	response := fiber.Map{
		"result":  product,
		"message": "Product was Inserted successfully.",
	}

	return c.Status(201).JSON(response)

}

func UpdateProduct(c *fiber.Ctx) error {

	id := c.Params("id")
	product := models.Product{}

	err := c.BodyParser(&product)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	db := database.DBConn
	productById := models.Product{}
	err = db.Find(&productById, id).Error
	if err != nil {
		// error handling...
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	if productById.ID == 0 {
		log.Print(productById)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Product not found"})
	}

	product.ID = productById.ID
	product.Id = productById.Id
	err = db.Model(&product).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		// error handling...
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	response := fiber.Map{
		"result":  product,
		"message": "Update Product Successfully.",
	}

	return c.JSON(response)

}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product := models.Product{}
	db := database.DBConn

	err := db.First(&product, id).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	err = db.Delete(&product, id).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	response := fiber.Map{
		"message": "Product Successfully deleted.",
	}

	return c.JSON(response)

}
