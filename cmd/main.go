package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	broker "ecommerce/broker"
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	"ecommerce/database"
	"ecommerce/pkg/config"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"

	logger "ecommerce/pkg/logger"

	validatorv10 "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"ecommerce/cmd/modules/admin"
	"ecommerce/cmd/modules/product"
	"ecommerce/cmd/modules/shop"
	"ecommerce/cmd/modules/transaction"
	"ecommerce/cmd/modules/user"
	"ecommerce/cmd/modules/warehouse"
)

func main() {
	swagger.New(swagger.Config{
		Title:        "Swagger API",
		DeepLinking:  false,
		DocExpansion: "none",
	})

	envConfig := config.SetupEnvFile()
	workerConfig := config.SetWorkerConfig()

	logger.InitializeLogger()

	fmt.Println("Loading database")
	db := database.InitDatabase(envConfig)

	fmt.Println("Migrate table")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	fmt.Println("Loading worker")
	rmq := broker.NewConnection(workerConfig)
	err := rmq.Connect()
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())

	app := fiber.New()
	routes.SetupRoutes(app)

	container := modules.SetContainerModules(app, db, validator, identifier, rmq, workerConfig)

	/* login route */
	user.InitModule(container)

	/* product route */
	product.InitModule(container)

	/* transaction route */
	transaction.InitModule(container)

	/* warehouse route */
	warehouse.InitModule(container)

	/* admin route */
	admin.InitModule(container)

	/* shop route */
	shop.InitModule(container)

	// Start the server in a separate goroutine
	go func() {
		if err := app.Listen(":5004"); err != nil {
			log.Fatalf("listen: %s", err)
		}
	}()

	// wait for termination signal and register database & http server clean-up operations
	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"http-server": func(ctx context.Context) error {
			return app.Shutdown()
		},
		// Add other cleanup operations here
	})

	<-wait
}

// operation is a clean up function on shutting down
type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
