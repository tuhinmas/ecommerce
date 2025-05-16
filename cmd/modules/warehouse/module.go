package warehouse

import (
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/warehouse"
	"ecommerce/internal/app/usecase/warehouse"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	warehouseRepository := repository.NewWarehouseRepository(container.Db)
	warehouseService := warehouse.NewWarehouseService(warehouseRepository, container.Validator)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	routes.WarehouseRouter(container.App, warehouseHandler)
}
