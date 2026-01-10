package response

import "time"

type CoffeeShopByIDResponse struct {
	ID        int                        `json:"id,omitempty"`
	Name      string                     `json:"name,omitempty"`
	Address   string                     `json:"address,omitempty"`
	Maps      string                     `json:"maps,omitempty"`
	Instagram string                     `json:"instagram,omitempty"`
	OpenClose string                     `json:"open_close,omitempty"`
	Category  string                     `json:"category,omitempty"`
	Facility  []FacilityResponse `json:"facility,omitempty"`
	CreatedBy *UserResponse               `json:"created_by,omitempty"`
	UpdatedBy *UserResponse               `json:"updated_by,omitempty"`
	UpdatedAt time.Time                  `json:"last_update,omitempty"`
	Images    []ImagesCoffeeShopResponse `json:"images"`
}

type CoffeeShopsResponse struct {
	ID        int                        `json:"id,omitempty"`
	Name      string                     `json:"name,omitempty"`
	Address   string                     `json:"address,omitempty"`
	OpenClose string                     `json:"open_close,omitempty"`
	Category  string                     `json:"category,omitempty"`
	Images    []ImagesCoffeeShopResponse `json:"images"`
}

type CreateCoffeeShopResponse struct {
	ID int `json:"id"`
}


