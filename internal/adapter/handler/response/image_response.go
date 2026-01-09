package response

type ImagesCoffeeShopResponse struct {
	Image     string `json:"image,omitempty"`
	IsPrimary bool   `json:"is_primary,omitempty"`
}