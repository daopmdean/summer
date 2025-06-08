package common

type Response struct {
	Status  ResponseStatusValue `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    interface{}         `json:"data,omitempty"`
	Paging  *PagingRes          `json:"paging,omitempty"`
	Errors  []*ErrRes           `json:"errors,omitempty"`
}

type PagingRes struct {
	TotalItems int64 `json:"totalItems,omitempty"`
	TotalPages int64 `json:"totalPages,omitempty"`
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
}

type ErrRes struct {
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
