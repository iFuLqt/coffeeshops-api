package response

type FacilityResponse struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}