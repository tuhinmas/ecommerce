package user

import (
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/user"
	"ecommerce/internal/app/usecase/user"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	repository := repository.NewUserRepository(container.Db)
	registerService := user.NewUserService(repository, container.Validator, container.Identifier)
	registerHandler := handler.NewUserHandler(registerService)
	routes.UserRouter(container.App, registerHandler)
}
