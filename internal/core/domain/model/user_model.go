package model

import "time"

type User struct {
	ID        int `gorm:"id"`
	Name      string `gorm:"name"`
	Email     string `gorm:"email"`
	Password  string `gorm:"password"`
	role      string `gorm:"role"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}