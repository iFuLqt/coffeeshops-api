package model

import "time"

type CoffeeShop struct {
	ID         int `gorm:"id"`
	Name       string `gorm:"name"`
	Address    string `gorm:"address"`
	Latitude   float64 `gorm:"latitude"`
	Longitude float64 `gorm:"longtitude"`
	OpenTime   string `gorm:"open_time"`
	CloseTime string `gorm:"close_time"`
	Parking bool `gorm:"parking"`
	PrayerRoom bool `gorm:"prayer_room"`
	Wifi bool `gorm:"wifi"`
	Gofood bool `gorm:"gofood"`
	Instagram string `gorm:"instagram"`
	CreatedByID int `gorm:"created_by_id"`
	UserCreate User `gorm:"foreignKey:CreatedByID"`
	UpdatedByID int `gorm:"updated_by_id"`
	UserUpdate User `gorm:"foreignKey:UpdatedByID"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	IsActive bool `gorm:"is_active"`
	CategoryID int `gorm:"category_id"`
	Category Category `gorm:"foreignKey:CategoryID"`
	Images []CoffeeShopImage `gorm:"foreignKey:CoffeeShopID"`
}