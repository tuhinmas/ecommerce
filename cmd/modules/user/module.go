package user

import (
	"ecommerce/cmd/routes"
	"ecommerce/database"
	repository "ecommerce/internal/app/repository/user"
	"ecommerce/internal/app/usecase/user"
	handler "ecommerce/internal/delivery"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func InitModule(app *fiber.App, db *database.Database, validator validator.Validator, identifier identifier.Identifier) {
	repository := repository.NewUserRepository(db)
	registerService := user.NewUserService(repository, validator, identifier)
	registerHandler := handler.NewUserHandler(registerService)
	routes.UserRouter(app, registerHandler)
}
