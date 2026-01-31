package models

import "time"

type Menu struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Code         string    `json:"code"`
	Route        string    `json:"route"`
	Icon         string    `json:"icon"`
	Level        int       `json:"level"`
	ReportID     int64     `json:"report_id"`
	ParentID     int64     `json:"parent_id"`
	DepartmentID int64     `json:"department_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Menu) Table() string {
	return "menus"
}
