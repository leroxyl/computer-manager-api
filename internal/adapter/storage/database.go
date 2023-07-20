package storage

import (
	"github.com/leroxyl/computer-manager-api/internal/config"
	"github.com/leroxyl/computer-manager-api/internal/entity"
	"github.com/leroxyl/computer-manager-api/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

const adminNotificationThreshold = 3 // TODO make admin notification threshold configurable

type notifyAdmin func(employeeAbbr string, computerCount int64)

type DatabaseManager struct {
	db *gorm.DB
	notifyAdmin
}

// Ensure DatabaseManager implements the ComputerManager interface
var _ usecase.ComputerManager = (*DatabaseManager)(nil)

// NewDatabaseManager initializes a new database connection and prepares the computer table
func NewDatabaseManager(conf config.DatabaseConfig, notify notifyAdmin) *DatabaseManager {
	// initialize db session
	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialize db session: %v", err)
	}

	// create or update computer table if necessary
	// TODO: make migration conditional
	err = db.AutoMigrate(&entity.Computer{})
	if err != nil {
		log.Fatalf("failed to auto-migrate database schema: %v", err)
	}

	return &DatabaseManager{
		db:          db,
		notifyAdmin: notify,
	}
}

func (dm *DatabaseManager) Create(computer entity.Computer) error {
	result := dm.db.Create(&computer)
	if result.Error != nil {
		return result.Error
	}

	go dm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (dm *DatabaseManager) Read(mac string) (entity.Computer, error) {
	computer := entity.Computer{
		MACAddr: mac,
	}

	result := dm.db.First(&computer)
	return computer, result.Error
}

func (dm *DatabaseManager) Update(computer entity.Computer) error {
	result := dm.db.Save(&computer)
	if result.Error != nil {
		return result.Error
	}

	go dm.checkComputerCount(computer.EmployeeAbbr)

	return nil
}

func (dm *DatabaseManager) Delete(mac string) error {
	computer := entity.Computer{
		MACAddr: mac,
	}

	result := dm.db.Delete(&computer)
	return result.Error
}

func (dm *DatabaseManager) ReadAll() (computers []entity.Computer, err error) {
	result := dm.db.Find(&computers)
	return computers, result.Error
}

func (dm *DatabaseManager) ReadAllForEmployee(employeeAbbr string) (computers []entity.Computer, err error) {
	result := dm.db.Where("employee_abbr = ?", employeeAbbr).Find(&computers)
	return computers, result.Error
}

func (dm *DatabaseManager) checkComputerCount(employeeAbbr string) {
	computerCount, err := dm.getComputerCountForEmployee(employeeAbbr)
	if err != nil {
		log.Errorf("failed to get computer count for employee %s from databse: %v", employeeAbbr, err)
		return
	}

	log.Infof("employee %s now has %d computers", employeeAbbr, computerCount)

	if computerCount >= adminNotificationThreshold {
		dm.notifyAdmin(employeeAbbr, computerCount)
	}
}

func (dm *DatabaseManager) getComputerCountForEmployee(employeeAbbr string) (count int64, err error) {
	err = dm.db.Model(&entity.Computer{}).Where("employee_abbr = ?", employeeAbbr).Count(&count).Error
	return count, err
}
