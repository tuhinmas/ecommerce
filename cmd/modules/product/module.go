package product

import (
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/product"
	service "ecommerce/internal/app/usecase/product"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	productRepository := repository.NewProductRepository(container.Db)
	productService := service.NewProductService(productRepository, container.Validator, container.Identifier)
	productHandler := handler.NewProductHandler(productService)
	routes.ProductRouter(container.App, productHandler)
}
