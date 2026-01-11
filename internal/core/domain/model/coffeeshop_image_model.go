package model

import "time"

type CoffeeShopImage struct {
	ID          int64 `gorm:"id"`
	CoffeeShopID int64 `gorm:"coffee_shop_id"`
	Image       string `gorm:"image"`
	IsPrimary   bool `gorm:"is_primary"`
	CreatedAt   time.Time `gorm:"created_at"`
}