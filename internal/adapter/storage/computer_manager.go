package storage

import (
	"os"

	"github.com/leroxyl/greenbone/internal/adapter/client"
	"github.com/leroxyl/greenbone/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

const (
	dsnEnv = "GREENBONE_POSTGRES_DSN"

	adminNotificationThreshold = 3 // TODO make admin notification threshold configurable
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
	// TODO: make migration conditional
	err = db.AutoMigrate(&entity.Computer{})
	if err != nil {
		log.Fatalf("failed to auto-migrate database schema: %v", err)
	}

	return &ComputerManager{
		db: db,
	}
}

func (cm *ComputerManager) Create(computer entity.Computer) error {
	result := cm.db.Create(&computer)
	if result.Error != nil {
		return result.Error
	}

	go cm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (cm *ComputerManager) Read(mac string) (entity.Computer, error) {
	computer := entity.Computer{
		MACAddr: mac,
	}

	result := cm.db.First(&computer)
	return computer, result.Error
}

func (cm *ComputerManager) Update(computer entity.Computer) error {
	result := cm.db.Save(&computer)
	if result.Error != nil {
		return result.Error
	}

	go cm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (cm *ComputerManager) Delete(mac string) error {
	computer := entity.Computer{
		MACAddr: mac,
	}

	result := cm.db.Delete(&computer)
	return result.Error
}

func (cm *ComputerManager) ReadAll() (computers []entity.Computer, err error) {
	result := cm.db.Find(&computers)
	return computers, result.Error
}

func (cm *ComputerManager) ReadAllForEmployee(employeeAbbr string) (computers []entity.Computer, err error) {
	result := cm.db.Where("employee_abbr = ?", employeeAbbr).Find(&computers)
	return computers, result.Error
}

func (cm *ComputerManager) checkComputerCount(employeeAbbr string) {
	computerCount, err := cm.getComputerCountForEmployee(employeeAbbr)
	if err != nil {
		log.Errorf("failed to get computer count for employee %s from databse: %v", employeeAbbr, err)
		return
	}

	log.Infof("employee %s now has %d computers", employeeAbbr, computerCount)

	if computerCount >= adminNotificationThreshold {
		client.NotifyAdmin(employeeAbbr, computerCount)
	}
}

func (cm *ComputerManager) getComputerCountForEmployee(employeeAbbr string) (count int64, err error) {
	err = cm.db.Model(&entity.Computer{}).Where("employee_abbr = ?", employeeAbbr).Count(&count).Error
	return count, err
}
