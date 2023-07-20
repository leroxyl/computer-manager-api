package usecase

import "github.com/leroxyl/computer-manager-api/internal/entity"

// ComputerManager interface represents the use cases for the Computer entity
type ComputerManager interface {
	Create(entity.Computer) error
	Read(mac string) (entity.Computer, error)
	Update(entity.Computer) error
	Delete(mac string) error
	ReadAll() ([]entity.Computer, error)
	ReadAllForEmployee(abbr string) ([]entity.Computer, error)
}
