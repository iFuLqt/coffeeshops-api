package response

type CategoryResponse struct {
	ID int `json:"id"`
	Category string `json:"category"`
	Slug string `json:"slug"`
	CreatedBy UserResponse `json:"created_by"`
}