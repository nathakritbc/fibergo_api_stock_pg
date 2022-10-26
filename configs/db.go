package configs

import (
	"fibergo_api_stock_pg/database"
	"fibergo_api_stock_pg/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() {
	var err error

	// ========== producttion mode =================
	// PSQL_HOST := Config("PSQL_HOST")
	// PSQL_DB_NAME := Config("PSQL_DB_NAME")
	// PSQL_USER := Config("PSQL_USER")
	// PSQL_PASS := Config("PSQL_PASS")
	// PSQL_PORT := Config("PSQL_PORT")

	// ================== dev mode =======================
	PSQL_HOST := Config("PSQL_HOST_DEV")
	PSQL_DB_NAME := Config("PSQL_DB_NAME_DEV")
	PSQL_USER := Config("PSQL_USER_DEV")
	PSQL_PASS := Config("PSQL_PASS_DEV")
	PSQL_PORT := Config("PSQL_PORT_DEV")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s",
		PSQL_HOST, PSQL_USER, PSQL_DB_NAME, PSQL_PORT, PSQL_PASS)

	log.Print("Connecting to PostgreSQL DB....")
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("connected")
	log.Println("running migrations")
	defer database.DBConn.AutoMigrate(&models.Product{}, &models.User{})

}
