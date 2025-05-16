package shop

import (
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/shop"
	"ecommerce/internal/app/usecase/shop"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	repository := repository.NewShopRepository(container.Db)
	shopService := shop.NewShopService(repository, container.Validator)
	shopHandler := handler.NewShopHandler(shopService)
	routes.ShopRouter(container.App, shopHandler)
}
