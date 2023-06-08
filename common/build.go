package common

func BuildMongoErr(msg string) *Response {
	return &Response{
		Status: ResponseStatus.Error,
		Error: &ErrorResponse{
			ErrorCode:    "MONGO_DB_ERROR",
			ErrorMessage: msg,
		},
	}
}

func BuildQueryNotFound(msg string) *Response {
	return &Response{
		Status:  ResponseStatus.NotFound,
		Message: msg,
		Error: &ErrorResponse{
			ErrorCode:    "QUERY_NOT_FOUND",
			ErrorMessage: msg,
		},
	}
}
