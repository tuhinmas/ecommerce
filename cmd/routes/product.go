package routes

import (
	"ecommerce/cmd/middleware"
	handler "ecommerce/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(app fiber.Router, productHandler handler.ProductHandler) {
	app.Get("/product", middleware.AuthUser, productHandler.ProductList)
	app.Post("/product", middleware.AuthAdmin, productHandler.CreateProduct)
}
