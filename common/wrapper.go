package common

func UnauthenticatedRes() *Response {
	return &Response{
		Status: ResponseStatus.Unauthorized,
		Errors: []*ErrRes{
			{
				ErrCode: "UNAUTHENTICATED",
				ErrMsg:  "Unauthenticated",
			},
		},
	}
}

func UnauthorizedRes() *Response {
	return &Response{
		Status: ResponseStatus.Unauthorized,
		Errors: []*ErrRes{
			{
				ErrCode: "UNAUTHORIZED",
				ErrMsg:  "Unauthorized",
			},
		},
	}
}

func InvalidRes(msg string) *Response {
	return &Response{
		Status: ResponseStatus.Invalid,
		Errors: []*ErrRes{
			{
				ErrCode: "INVALID",
				ErrMsg:  msg,
			},
		},
	}
}
