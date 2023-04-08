package common

type Response struct {
	Status  ResponseStatusValue `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    interface{}         `json:"data,omitempty"`
	Paging  *PagingResponse     `json:"paging,omitempty"`
	Error   *ErrorResponse      `json:"error,omitempty"`
}

type PagingResponse struct {
	Total int64 `json:"total,omitempty"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
