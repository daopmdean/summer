package common

func BuildMongoErr(msg string) *Response {
	return &Response{
		Status: ResponseStatus.Error,
		Error: &ErrorResponse{
			ErrorMessage: msg,
			ErrorCode:    "ERR_MONGO_DB",
		},
	}
}
