package controllers

import (
	_ "fibergo_api_stock_pg/helpers"

	"mime/multipart"
	"os"
	"strings"

	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
)

func checkImageType(imageType string, c *fiber.Ctx) error {

	runes := []rune(imageType)
	imgType := strings.TrimSpace(string(runes[0:5]))
	if imgType != "image" {
		return fiber.NewError(fiber.StatusBadRequest, "error not found image type")
	}
	return nil
}

func uploadfile(file *multipart.FileHeader, newFileName string, c *fiber.Ctx) error {
	c.SaveFile(file, fmt.Sprintf("./images/%s", newFileName))
	return nil
}

func UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
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

	uuid := guuid.New().String()
	newFileName := uuid + file.Filename
	uploadfile(file, newFileName, c)

	return c.Status(200).JSON(fiber.Map{"fileName": newFileName, "message": "upload files successfully"})

}

func RemoveImage(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	filePath := "./images/" + fileName
	_, e := os.Stat(filePath)
	if e != nil {
		fmt.Printf("File does not exist")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}

	err := os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}

	messageResponse := "remove files " + fileName + " successfully"

	return c.Status(200).JSON(fiber.Map{"message": messageResponse})
}
