package routes

import (
	"ecommerce/cmd/middleware"
	handler "ecommerce/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

func ShopRouter(app fiber.Router, shopHandler handler.ShopHandler) {
	app.Post("/shop", middleware.AuthAdmin, shopHandler.CreateShop)
}
