package model

import "time"

type Category struct {
	ID        int `gorm:"id"`
	Name      string `gorm:"name"`
	Slug      string `gorm:"slug"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	CreatedByID int `gorm:"created_by_id"`
	User User `gorm:"foreignKey:CreatedByID"`
}