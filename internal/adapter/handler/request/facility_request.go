package request 

type FacilityCoffeeShopRequest struct {
	Code []string `json:"facility_code" validate:"required,min=1,dive,required"`
}