package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

const (
	dsnEnv = "GREENBONE_POSTGRES_DSN"

	adminNotificationThreshold = 2 // TODO externalize admin notification threshold
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
	if result.Error != nil {
		return result.Error
	}

	go cm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (cm *ComputerManager) Read(computer *Computer) error {
	result := cm.db.First(computer)
	return result.Error
}

func (cm *ComputerManager) Update(computer *Computer) error {
	result := cm.db.Save(computer)
	if result.Error != nil {
		return result.Error
	}

	go cm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (cm *ComputerManager) Delete(computer *Computer) error {
	result := cm.db.Delete(computer)
	return result.Error
}

func (cm *ComputerManager) ReadAll(computers *[]Computer) error {
	result := cm.db.Find(computers)
	return result.Error
}

func (cm *ComputerManager) checkComputerCount(employeeAbbr string) {
	computerCount, err := cm.getComputerCountForEmployee(employeeAbbr)
	if err != nil {
		log.Errorf("failed to get computer count for employee %s from databse: %v", employeeAbbr, err)
		return
	}

	log.Infof("employee %s now has %d computers", employeeAbbr, computerCount)

	if computerCount > adminNotificationThreshold {
		notifyAdmin(employeeAbbr, computerCount)
	}
}

func (cm *ComputerManager) getComputerCountForEmployee(employeeAbbr string) (count int64, err error) {
	err = cm.db.Model(&Computer{}).Where("employee_abbr = ?", employeeAbbr).Count(&count).Error
	return count, err
}
