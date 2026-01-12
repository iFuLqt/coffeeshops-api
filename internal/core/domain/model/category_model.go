package model

import "time"

type Category struct {
	ID        int64 `gorm:"id"`
	Name      string `gorm:"name"`
	Slug      string `gorm:"slug"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	CreatedByID int64 `gorm:"created_by_id"`
	User User `gorm:"foreignKey:CreatedByID"`
	CoffeeShops []CoffeeShop `gorm:"foreignKey:CategoryID"`
}