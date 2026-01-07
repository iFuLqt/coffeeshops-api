package response

type CategoryResponse struct {
	ID int `json:"id"`
	Category string `json:"category"`
	Slug string `json:"slug"`
	CreatedByName string `json:"created_by_name"`
}