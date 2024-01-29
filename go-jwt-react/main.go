package main

import (
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/database"
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":8000")
}
