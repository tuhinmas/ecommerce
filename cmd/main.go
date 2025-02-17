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
	"ecommerce/cmd/routes"
	"ecommerce/database"
	adminRepository "ecommerce/internal/app/repository/admin"
	productRepository "ecommerce/internal/app/repository/product"
	shopRepository "ecommerce/internal/app/repository/shop"
	transactionRepository "ecommerce/internal/app/repository/transaction"
	repository "ecommerce/internal/app/repository/user"
	warehouseRepository "ecommerce/internal/app/repository/warehouse"
	"ecommerce/internal/app/usecase/admin"
	"ecommerce/internal/app/usecase/product"
	"ecommerce/internal/app/usecase/shop"
	"ecommerce/internal/app/usecase/transaction"
	"ecommerce/internal/app/usecase/user"
	"ecommerce/internal/app/usecase/warehouse"
	"ecommerce/internal/app/worker"
	handler "ecommerce/internal/delivery"
	"ecommerce/pkg/config"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"

	logger "ecommerce/pkg/logger"

	validatorv10 "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
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

	fmt.Println("Loading worker")
	rmq := broker.NewConnection(workerConfig)
	err := rmq.Connect()
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())

	repository := repository.NewUserRepository(db)
	adminRepository := adminRepository.NewAdminRepository(db)
	productRepository := productRepository.NewProductRepository(db)
	transactionRepository := transactionRepository.NewTransactionRepository(db)
	warehouseRepository := warehouseRepository.NewWarehouseRepository(db)
	shopRepository := shopRepository.NewShopRepository(db)

	queueService := worker.NewQueueService(rmq, workerConfig, transactionRepository)
	registerService := user.NewUserService(repository, validator, identifier)
	transactionService := transaction.NewTransactionService(transactionRepository, validator, identifier, queueService)
	productService := product.NewProductService(productRepository, validator, identifier)
	warehouseService := warehouse.NewWarehouseService(warehouseRepository, validator)
	adminService := admin.NewAdminService(adminRepository, validator, identifier)
	shopService := shop.NewShopService(shopRepository, validator)

	registerHandler := handler.NewUserHandler(registerService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	productHandler := handler.NewProductHandler(productService)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	adminHandler := handler.NewAdminHandler(adminService)
	shopHandler := handler.NewShopHandler(shopService)

	app := fiber.New()
	routes.SetupRoutes(app)
	routes.UserRouter(app, registerHandler)
	routes.TransactionRouter(app, transactionHandler)
	routes.ProductRouter(app, productHandler)
	routes.WarehouseRouter(app, warehouseHandler)
	routes.AdminRouter(app, adminHandler)
	routes.ShopRouter(app, shopHandler)

	go queueService.ConsumeData(context.Background(), workerConfig.StockReversalQueue)

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
