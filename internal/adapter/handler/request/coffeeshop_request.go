package request

type CoffeeShopRequest struct {
	CoffeeShop string  `json:"coffee_shop" validate:"required"`
	Address    string  `json:"address" validate:"required"`
	Latitude   float64 `json:"latitude" validate:"required"`
	Longitude  float64 `json:"longitude" validate:"required"`
	OpenTime   string  `json:"open_time" validate:"required"`
	CloseTime  string  `json:"close_time" validate:"required"`
	Parking    bool    `json:"parking" validate:"required"`
	PrayerRoom bool    `json:"prayer_room" validate:"required"`
	Wifi       bool    `json:"wifi" validate:"required"`
	Gofood     bool    `json:"gofood" validate:"required"`
	Instagram  string  `json:"instagram" validate:"required"`
	CategoryID int     `json:"category_id" validate:"required"`
}
