package entity

type CoffeeShopFacilityEntity struct {
	ID           int64
	CoffeeShopID int64
	FacilityID   int64
	Facility     FacilityEntity
}
