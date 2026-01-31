package models

import "time"

type Department struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Code      string    `json:"code" gorm:"not null"`
	Desc      string    `json:"desc" gorm:"varchar(255); not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Department) Table() string {
	return "departments"
}
