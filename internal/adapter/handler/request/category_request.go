package request

type CategoryRequest struct {
	Category string `json:"category" validate:"required"`
}