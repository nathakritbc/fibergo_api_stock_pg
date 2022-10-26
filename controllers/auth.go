package controllers

import (
	"log"
	"time"

	"fibergo_api_stock_pg/configs"
	"fibergo_api_stock_pg/database"
	"fibergo_api_stock_pg/helpers"
	"fibergo_api_stock_pg/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
)

func Register(c *fiber.Ctx) error {
	user := models.User{}

	err := c.BodyParser(&user)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	errors := helpers.ValidateStructUser(user)
	if errors != nil {
		log.Print(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors})

	}

	user.U_password, _ = hashPassword(user.U_password)

	db := database.DBConn
	err = db.Create(&user).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User was registered successfully!"})
}

func Login(c *fiber.Ctx) error {

	user := models.User{}
	err := c.BodyParser(&user)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	errors := helpers.ValidateStructUser(user)
	if errors != nil {
		log.Print(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors})
	}

	passwordInput := user.U_password
	user.U_password, _ = hashPassword(user.U_password)
	db := database.DBConn
	err = db.First(&user, "u_email = ?", user.U_email).Error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "User Not found."})
	}

	if !checkPasswordHash(passwordInput, user.U_password) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Invalid Password!"})

	}

	// Create the Claims
	claims := jwt.MapClaims{
		"u_id":    user.ID,
		"u_email": user.U_email,
		"admin":   false,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	var secretKey = configs.Config("SECRET_KEY")
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	//	create response
	response := fiber.Map{"result": fiber.Map{
		"token":   t,
		"u_email": user.U_email,
		"u_id":    user.ID,
	}}

	return c.JSON(response)

}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	return string(bytes), err
}

func DecodeUser(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userMap := fiber.Map{"user": fiber.Map{
		"u_id":    claims["u_id"],
		"u_email": claims["u_email"],
	}}

	return c.JSON(userMap)

}
