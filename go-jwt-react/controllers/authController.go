package controllers

import (
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/database"
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(404).SendString("Could not hash password")
	}
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DB.Create(&user)
	return c.Status(200).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(user)
	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadGateway)
		return c.JSON(&fiber.Map{
			"message": "Incorrect Password",
		})
	}

	return c.Status(200).JSON(user)
}
