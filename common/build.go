package common

func BuildMongoErr(msg string) *Response {
	return &Response{
		Status: ResponseStatus.Error,
		Errors: []*ErrorRes{
			{
				ErrCode: "MONGO_DB_ERROR",
				ErrMsg:  msg,
			},
		},
	}
}

func BuildQueryNotFound(msg string) *Response {
	return &Response{
		Status:  ResponseStatus.NotFound,
		Message: msg,
		Errors: []*ErrorRes{
			{
				ErrCode: "QUERY_NOT_FOUND",
				ErrMsg:  msg,
			},
		},
	}
}
