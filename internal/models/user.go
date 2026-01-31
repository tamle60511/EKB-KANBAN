package models

import "time"

type User struct {
	ID           int64       `json:"id"`
	Username     string      `json:"username"`
	Password     string      `json:"password"`
	FullName     string      `json:"full_name"`
	Email        string      `json:"email,omitempty"`
	DepartmentID int64       `json:"department_id,omitempty"`
	Department   *Department `json:"department,omitempty"`
	IsActive     bool        `json:"is_active" gorm:"default:true"`
	Role         string      `json:"role" gorm:"default:'user'"`
	LastLogin    time.Time   `json:"last_login,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (User) Table() string {
	return "users"
}
