package models

import "time"

type Report struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	ReportType     string    `gorm:"type:varchar(100);not null" json:"report_type"`
	ReportName     string    `gorm:"type:varchar(255);not null" json:"report_name"`
	DepartmentID   string    `josn:"department_id" gorm:"not null"`
	QueryStatement string    `gorm:"type:text;not null" json:"query_statement"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Report) Table() string {
	return "reports"
}

type ReportColumn struct {
	ID       int64  `json:"id" gorm:"primaryKey" `
	ReportID int64  `json:"report_id" gorm:"not null"`
	Title    string `json:"title" gorm:"not null"`
	Code     string `json:"code" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
	Num      int64  `json:"num" gorm:"not null"`
}

func (ReportColumn) Table() string {
	return "report_columns"
}
