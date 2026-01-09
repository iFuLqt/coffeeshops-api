package model

type CoffeeShopFacility struct {
	ID           int `gorm:"id"`
	CoffeeShopID int `gorm:"coffee_shop_id"`
	FacilityID   int `gorm:"facility_id"`
	Facility Facility `gorm:"foreignKey:FacilityID"`
}
