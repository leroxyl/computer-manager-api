package database

type Computer struct {
	MACAddr      string `json:"macAddr" gorm:"primaryKey;check:mac_addr_required,mac_addr <> ''"`
	ComputerName string `json:"computerName" gorm:"not null;check:computer_name_required,computer_name <> ''"`
	IPAddr       string `json:"ipAddr" gorm:"not null;check:ip_addr_required,ip_addr <> ''"`
	EmployeeAbbr string `json:"employeeAbbr"`
	Description  string `json:"description"`
}
