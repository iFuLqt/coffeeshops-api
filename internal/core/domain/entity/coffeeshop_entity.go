package entity

type CoffeeShopEntity struct {
	Name       string
	Address    string
	Latitude   float64
	Longitude  float64
	OpenTime   string
	CloseTime  string
	Parking    bool
	PrayerRoom bool
	Wifi       bool
	Gofood     bool
	Instagram  string
	UserCreate UserEntity
	UserUpdate UserEntity
	Category   CategoryEntity
	Image      []ImageEntity
	IsActive bool
}
