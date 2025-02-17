package database

import (
	"fmt"

	"ecommerce/pkg/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

func InitDatabase(config *config.Config) *Database {
	dsn := getDataSourceName(config)
	database, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}

	return &Database{
		database,
	}
}

func getDataSourceName(config *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
	)
}
