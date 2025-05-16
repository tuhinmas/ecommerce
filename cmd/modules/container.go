package modules

import (
	"ecommerce/broker"
	"ecommerce/database"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type Container struct {
	App          *fiber.App
	Db           *database.Database
	Validator    validator.Validator
	Identifier   identifier.Identifier
	Rmq          broker.RabbitMQ
	WorkerConfig broker.RabbitmqConfig
}

func SetContainerModules(
	app *fiber.App,
	database *database.Database,
	validator validator.Validator,
	identifier identifier.Identifier,
	rmq broker.RabbitMQ,
	workerConfig broker.RabbitmqConfig,
) Container {
	return Container{
		App:          app,
		Db:           database,
		Validator:    validator,
		Identifier:   identifier,
		Rmq:          rmq,
		WorkerConfig: workerConfig,
	}
}
