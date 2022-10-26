package controllers

import (
	"fibergo_api_stock_pg/database"
	"fibergo_api_stock_pg/helpers"
	_ "fibergo_api_stock_pg/helpers"
	"fibergo_api_stock_pg/models"
	"mime/multipart"
	"strings"

	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
)

func removeImage(removePath string, c *fiber.Ctx) error {
	e := os.Remove(removePath)
	if e != nil {
		log.Fatal(e)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	return nil
}

func checkImageType(imageType string, c *fiber.Ctx) error {

	runes := []rune(imageType)
	imgType := strings.TrimSpace(string(runes[0:5]))
	if imgType != "image" {
		return fiber.NewError(fiber.StatusBadRequest, "error not found image type")
	}
	return nil
}

func uploadImage(file *multipart.FileHeader, newFileName string, c *fiber.Ctx) error {
	c.SaveFile(file, fmt.Sprintf("./images/%s", newFileName))
	return nil
}

func GetProductsAll(c *fiber.Ctx) error {
	product := []models.Product{}
	db := database.DBConn
	err := db.Find(&product).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	if product.ID == 0 {
		log.Print(product)
		return c.JSON(fiber.Map{})
	}

	return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {

	product := models.Product{}

	file, err := c.FormFile("p_image")
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	imageType := file.Header["Content-Type"][0]
	err = checkImageType(imageType, c)
	if err != nil {
		file = nil
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	err = c.BodyParser(&product)
	if err != nil {
		file = nil
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	errors := helpers.ValidateStructProduct(product)
	if errors != nil {
		log.Print(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors})
	}

	uuid := guuid.New().String()
	newFileName := uuid + file.Filename
	uploadImage(file, newFileName, c)

	product.P_Image = newFileName
	db := database.DBConn
	err = db.Create(&product).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	errors := helpers.ValidateStructProduct(product)
	if errors != nil {
		log.Print(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors})
	}

	db := database.DBConn
	productById := models.Product{}
	err = db.Find(&productById, id).Error
	if err != nil {
		// error handling...
		log.Print(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	if productById.ID == 0 {
		log.Print(productById)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	file, errImgForm := c.FormFile("p_image")

	newFileName := ""
	if errImgForm != nil {
		log.Print(err)
		product.P_Image = productById.P_Image
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	} else {
		uuid := guuid.New().String()
		newFileName = uuid + file.Filename
		product.P_Image = newFileName
		uploadImage(file, newFileName, c)

	}

	// product.P_Image = newFileName
	product.ID = productById.ID
	product.Id = productById.Id
	err = db.Model(&product).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		// error handling...
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	if errImgForm == nil {
		removePath := "./images/" + productById.P_Image
		removeImage(removePath, c)
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	err = db.Delete(&product, id).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	removePath := "./images/" + product.P_Image
	removeImage(removePath, c)

	response := fiber.Map{
		"message": "Product Successfully deleted.",
	}

	return c.JSON(response)

}
