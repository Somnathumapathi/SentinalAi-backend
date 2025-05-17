package db

import (
	"fmt"
	"log"
	"os"
	"sentinal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitGORM() error {
	// Get database connection details from environment variables
	host := os.Getenv("SUPABASE_DB_HOST")
	user := os.Getenv("SUPABASE_DB_USER")
	password := os.Getenv("SUPABASE_DB_PASSWORD")
	dbname := os.Getenv("SUPABASE_DB_NAME")
	port := os.Getenv("SUPABASE_DB_PORT")

	// Create DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		host, user, password, dbname, port)

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Enable UUID extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.GitHubInstallation{},
		&models.GitHubRepository{},
		&models.AWSIntegration{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Connected to database with GORM!")
	return nil
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return DB
}
