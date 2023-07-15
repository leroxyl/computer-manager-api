package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dsnEnv = "GREENBONE_POSTGRES_DSN"
)

type ComputerManager struct {
	db *gorm.DB
}

// NewComputerManager initializes a new database connection and prepares the computer table
func NewComputerManager() *ComputerManager {
	// load DSN from environment
	dsn := os.Getenv(dsnEnv)
	if dsn == "" {
		log.Fatalf("environment variable %s not set, please provide DSN for postgres database instance", dsnEnv)
	}

	// initialize db session
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialize db session: %v", err)
	}

	// create or update computer table if necessary
	err = db.AutoMigrate(&Computer{})
	if err != nil {
		log.Fatalf("failed to auto-migrate database schema: %v", err)
	}

	return &ComputerManager{
		db: db,
	}
}

func (cm *ComputerManager) Create(computer *Computer) error {
	result := cm.db.Create(computer)
	return result.Error
}
