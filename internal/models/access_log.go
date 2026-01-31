package models

import "time"

// AccessLog represents a log of user access to operations
type AccessLog struct {
	ID           int         `json:"id"`
	UserID       int64       `json:"user_id"`
	DepartmentID int64       `json:"department_id"`
	OperationID  int         `json:"operation_id"`
	AccessTime   time.Time   `json:"access_time"`
	ReportID     int64       `json:"report_id,omitempty"`
	IPAddress    string      `json:"ip_address,omitempty"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Operation    *Operation  `gorm:"foreignKey:OperationID" json:"operation,omitempty"`
}

func (a *AccessLog) TableName() string {
	return "access_logs"
}
