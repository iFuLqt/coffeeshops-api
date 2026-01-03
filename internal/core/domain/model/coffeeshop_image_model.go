package model

import "time"

type CoffeeShopImage struct {
	ID          string `gorm:"id"`
	CoffeeShopID int `gorm:"coffee_shop_id"`
	Image       string `gorm:"image"`
	IsPrimary   bool `gorm:"is_primary"`
	CreatedAt   time.Time `gorm:"created_at"`
}