package entity

import "time"

type CoffeeShopEntity struct {
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
	OpenTime  time.Time
	CloseTime time.Time
	Parking bool
	PrayerRoom bool
	Wifi bool
	Gofood bool
	Instagram string
	UserCreate UserEntity
	UserUpdate UserEntity
	Category CategoryEntity
}