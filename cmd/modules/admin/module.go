package admin

import (
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/admin"
	"ecommerce/internal/app/usecase/admin"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	repository := repository.NewAdminRepository(container.Db)
	adminService := admin.NewAdminService(repository, container.Validator, container.Identifier)
	adminHandler := handler.NewAdminHandler(adminService)
	routes.AdminRouter(container.App, adminHandler)
}
