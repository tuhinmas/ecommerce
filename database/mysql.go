package database

import (
	"fmt"
	"time"

	"ecommerce/pkg/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	*sqlx.DB
}

func InitDatabase(config *config.Config) *Database {
	dsn := getDataSourceName(config)
	database, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		// panic("failed to connect database")
		panic(err.Error())
	}

	// Konfigurasi connection pool
	database.SetMaxOpenConns(100)                // Maksimal koneksi yang dibuka / dibuat
	database.SetMaxIdleConns(25)                 // Maksimal koneksi idle (connection yang menunggu di pool)
	database.SetConnMaxLifetime(5 * time.Minute) // Maksimal umur koneksi
	database.SetConnMaxIdleTime(2 * time.Minute) // Maksimal waktu koneksi idle

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
func RunMigrations(db *Database) error {
	sqlDB := db.DB

	driver, err := mysql.WithInstance(sqlDB.DB, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
