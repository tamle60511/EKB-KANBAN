package dto

import "time"

type AccessLog struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	DepartmentID int64     `json:"department_id"`
	OperationID  int64     `json:"operation_id"`
	AccessTime   time.Time `json:"access_time"`
	ReportID     int64     `json:"report_id"`
	IPAddress    string    `json:"ip_address"`
	Status       string    `json:"status"`
}

type AccessLogResponse struct {
	ID             int       `json:"id"`
	UserID         int64     `json:"user_id"`
	Username       string    `json:"username,omitempty"`
	DepartmentID   int64     `json:"department_id"`
	DepartmentName string    `json:"department_name,omitempty"`
	OperationID    int       `json:"operation_id"`
	OperationName  string    `json:"operation_name,omitempty"`
	OperationCode  string    `json:"operation_code,omitempty"`
	AccessTime     time.Time `json:"access_time"`
	ReportID       int64     `json:"report_id,omitempty"`
	IPAddress      string    `json:"ip_address,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type OperationResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Desc        string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessCount int64     `json:"access_count,omitempty"`
}

type DailyAccessCount struct {
	Date  string `gorm:"column:date"`
	Count int64  `gorm:"column:count"`
}
