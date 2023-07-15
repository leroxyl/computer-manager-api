package database

type Computer struct {
	MACAddr      string `json:"macAddr" gorm:"primaryKey"`
	ComputerName string `json:"computerName" gorm:"not null"`
	IPAddr       string `json:"ipAddr" gorm:"not null"`
	EmployeeAbbr string `json:"employeeAbbr"`
	Description  string `json:"description"`
}
