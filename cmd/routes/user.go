package routes

import (
	handler "ecommerce/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, signupHandler handler.UserHandler) {
	app.Post("/register", signupHandler.Signup)
	app.Post("/login", signupHandler.Login)
}

func AdminRouter(app fiber.Router, adminHandler handler.AdminHandler) {
	app.Post("/admin/login", adminHandler.Login)
}
