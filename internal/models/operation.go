package models

import "time"

type Operation struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	AccessLogs  []AccessLog `gorm:"foreignKey:OperationID" json:"access_logs,omitempty"`
}

func (o *Operation) TableName() string {
	return "operations"
}
