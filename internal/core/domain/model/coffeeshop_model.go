package model

import "time"

type CoffeeShop struct {
	ID                 int64                `gorm:"id"`
	Name               string               `gorm:"name"`
	Slug               string               `gorm:"slug"`
	Address            string               `gorm:"address"`
	Latitude           float64              `gorm:"latitude"`
	Longitude          float64              `gorm:"longtitude"`
	OpenTime           string               `gorm:"open_time"`
	CloseTime          string               `gorm:"close_time"`
	Instagram          string               `gorm:"instagram"`
	CreatedByID        int64                `gorm:"created_by_id"`
	UserCreate         User                 `gorm:"foreignKey:CreatedByID"`
	UpdatedByID        int64                `gorm:"updated_by_id"`
	UserUpdate         User                 `gorm:"foreignKey:UpdatedByID"`
	CreatedAt          time.Time            `gorm:"created_at"`
	UpdatedAt          time.Time            `gorm:"updated_at"`
	IsActive           *bool                 `gorm:"is_active"`
	CategoryID         int64                `gorm:"category_id"`
	Category           Category             `gorm:"foreignKey:CategoryID"`
	Images             []CoffeeShopImage    `gorm:"foreignKey:CoffeeShopID"`
	CoffeeShopFacility []CoffeeShopFacility `gorm:"foreignKey:CoffeeShopID"`
}
