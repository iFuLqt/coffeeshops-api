package entity

import "time"

type CoffeeShopEntity struct {
	ID         int64
	Name       string
	Slug       string
	Address    string
	Latitude   float64
	Longitude  float64
	OpenTime   string
	CloseTime  string
	Instagram  string
	UserCreate UserEntity
	UserUpdate UserEntity
	Category   CategoryEntity
	Image      []ImageEntity
	Facility   []FacilityEntity
	IsActive   *bool
	UpdatedAt  time.Time
}
