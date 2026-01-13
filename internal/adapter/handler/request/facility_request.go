package request 

type FacilityCoffeeShopRequest struct {
	FacilityCode []string `json:"facility_code" validate:"required"`
}

type FacilityRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}