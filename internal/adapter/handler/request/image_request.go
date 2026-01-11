package request

type ImageRequest struct {
	Image []int64 `json:"image_ids" validate:"required,min=1,dive,gt=0"`
}