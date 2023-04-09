package common

func ResponseError(msg string) *Response {
	return &Response{
		Status:  ResponseStatus.Error,
		Message: msg,
	}
}
