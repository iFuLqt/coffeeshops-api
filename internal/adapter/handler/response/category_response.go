package response

type CategoryResponse struct {
	ID int64 `json:"id"`
	Category string `json:"category"`
	Slug string `json:"slug"`
	CreatedBy UserResponse `json:"created_by"`
}