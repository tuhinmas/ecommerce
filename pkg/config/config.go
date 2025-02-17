package config

import (
	"ecommerce/broker"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseUser     string `envconfig:"DATABASE_USER" default:"root"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" default:"opklnm123"`
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"mysql-local"`
	DatabaseName     string `envconfig:"DATABASE_NAME" default:"dating"`
	DatabasePort     string `envconfig:"DATABASE_PORT" default:"5432"`
}

func SetupEnvFile() *Config {
	envConfig := &Config{}
	_ = godotenv.Load()
	err := envconfig.Process("", envConfig)
	if err != nil {
		log.Fatal(nil, "Fatal error ", err)
	}

	return envConfig
}

func SetWorkerConfig() broker.RabbitmqConfig {
	port, err := strconv.Atoi(GetEnv("WORKER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return broker.RabbitmqConfig{
		Host:                    GetEnv("WORKER_HOST"),
		Username:                GetEnv("WORKER_USERNAME"),
		Password:                GetEnv("WORKER_PASSWORD"),
		Port:                    port,
		StockReversalExchange:   GetEnv("WORKER_STOCK_REVERSAL_EXCHANGE"),
		StockReversalQueue:      GetEnv("WORKER_STOCK_REVERSAL_QUEUE"),
		StockReversalRoutingKey: GetEnv("WORKER_STOCK_REVERSAL_ROUTING_KEY"),
		PendingPaymentQueue:     GetEnv("WORKER_PENDING_PAYMENT_QUEUE"),
	}
}

// GetEnv
func GetEnv(key string) string {
	value := os.Getenv(key)
	return value
}
