package routes

import (
	"ecommerce/cmd/middleware"
	handler "ecommerce/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

func WarehouseRouter(app fiber.Router, warehouseHandler handler.WarehouseHandler) {
	app.Post("/warehouse", middleware.AuthAdmin, warehouseHandler.CreateWarehouse)
	app.Post("/warehouse/status", middleware.AuthAdmin, warehouseHandler.SetStatusWarehouse)
	app.Post("/warehouse/stock", middleware.AuthAdmin, warehouseHandler.CreateStock)
	app.Put("/warehouse/stock/:id", middleware.AuthAdmin, warehouseHandler.UpdateStock)
	app.Post("/warehouse/stock-transfer", middleware.AuthAdmin, warehouseHandler.TransferStock)
}
