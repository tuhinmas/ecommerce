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
	DatabaseUser     string `envconfig:"DB_USERNAME" default:"root"`
	DatabasePassword string `envconfig:"DB_PASSWORD" default:"root"`
	DatabaseHost     string `envconfig:"DB_HOST" default:"mysql"`
	DatabaseName     string `envconfig:"DB_DATABASE" default:"e-commerce"`
	DatabasePort     string `envconfig:"DB_PORT" default:"3306"`
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
