package response

type FacilityResponse struct {
	ID int64 `json:"id,omitempty"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}