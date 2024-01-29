package routes

import (
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Group("/api")
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.Logout)
}
