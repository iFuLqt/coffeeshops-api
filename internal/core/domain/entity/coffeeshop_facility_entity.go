package entity

type CoffeeShopFacilityEntity struct {
	ID           int
	CoffeeShopID int
	FacilityID   int
	Facility     FacilityEntity
}
