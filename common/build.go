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
