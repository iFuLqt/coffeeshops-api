package model

type CoffeeShopFacility struct {
	ID           int64 `gorm:"id"`
	CoffeeShopID int64 `gorm:"coffee_shop_id"`
	CoffeeShop CoffeeShop `gorm:"foreignKey:CoffeeShopID"`
	FacilityID   int64 `gorm:"facility_id"`
	Facility Facility `gorm:"foreignKey:FacilityID"`
}

func (CoffeeShopFacility) TableName() string {
	return "coffee_shop_facility"
}
