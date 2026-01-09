package model

type CoffeeShopFacility struct {
	ID           int `gorm:"id"`
	CoffeeShopID int `gorm:"coffee_shop_id"`
	CoffeeShop CoffeeShop `gorm:"foreignKey:CoffeeShopID"`
	FacilityID   int `gorm:"facility_id"`
	Facility Facility `gorm:"foreignKey:FacilityID"`
}

func (CoffeeShopFacility) TableName() string {
	return "coffee_shop_facility"
}
