package routes

import (
	"ecommerce/cmd/middleware"
	handler "ecommerce/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

func TransactionRouter(app fiber.Router, transactionHandler handler.TransactionHandler) {
	app.Post("/transaction", middleware.AuthUser, transactionHandler.Order)
}
