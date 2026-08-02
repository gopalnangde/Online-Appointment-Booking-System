package database

import (
	"fmt"
	"log"

	"backend/config"
	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection instance used across the application
var DB *gorm.DB

// Connect establishes a connection to the MySQL database and runs auto-migrations.
// It uses the DSN (Data Source Name) constructed from the provided Config.
func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred",
	cfg.DBUser,
	cfg.DBPassword,
	cfg.DBHost,
	cfg.DBPort,
	cfg.DBName,
)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Auto-migrate all models to create/update the database tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.ServiceProviderProfile{},
		&models.Appointment{},
		&models.Review{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
}

// GetDB returns the database connection instance.
// This is preferred over directly accessing the DB variable for testability.
func GetDB() *gorm.DB {
	return DB
}
