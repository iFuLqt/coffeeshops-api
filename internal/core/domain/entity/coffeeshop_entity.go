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

type QueryString struct {
	Limit int64
	Page int64
	OrderBy string
	OrderType string
	Search string
	CategoryID int64
	Status bool
}
