package response

type Meta struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Errors interface{} `json:"errors"`
}

type DefaultErrorResponse struct {
	Meta Meta `json:"meta"`
}

type DefaultSuccessResponse struct {
	Meta Meta `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}