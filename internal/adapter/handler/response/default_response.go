package response

type Meta struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type DefaultErrorResponse struct {
	Meta Meta `json:"meta"`
}

type DefaultSuccessResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	TotalRecords int64 `json:"total_records"`
	Page         int64 `json:"page"`
	PerPage      int64 `json:"per_page"`
	TotalPages   int64 `json:"total_pages"`
}
