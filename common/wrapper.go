package common

func UnauthenticatedRes() *Response {
	return BuildUnauthorizedRes("UNAUTHENTICATED", "Unauthenticated")
}

func UnauthorizedRes() *Response {
	return BuildUnauthorizedRes("UNAUTHORIZED", "Unauthorized")
}

func BuildUnauthorizedRes(errCode, errMsg string) *Response {
	return &Response{
		Status: ResponseStatus.Unauthorized,
		Errors: []*ErrRes{
			{
				ErrCode: errCode,
				ErrMsg:  errMsg,
			},
		},
	}
}

func InvalidRes(msg string) *Response {
	return BuildInvalidRes("INVALID", msg)
}

func BuildInvalidRes(errCode, errMsg string) *Response {
	return &Response{
		Status: ResponseStatus.Invalid,
		Errors: []*ErrRes{
			{
				ErrCode: errCode,
				ErrMsg:  errMsg,
			},
		},
	}
}

func BuildErrorRes(errCode, errMsg string) *Response {
	return &Response{
		Status: ResponseStatus.Error,
		Errors: []*ErrRes{
			{
				ErrCode: errCode,
				ErrMsg:  errMsg,
			},
		},
	}
}

func BuildErrorRes(errCode, errMsg string) *Response {
	return &Response{
		Status: ResponseStatus.NotFound,
		Errors: []*ErrRes{
			{
				ErrCode: errCode,
				ErrMsg:  errMsg,
			},
		},
	}
}

func BuildRes(status ResponseStatusValue, errCode, errMsg string) *Response {
	return &Response{
		Status: status,
		Errors: []*ErrRes{
			{
				ErrCode: errCode,
				ErrMsg:  errMsg,
			},
		},
	}
}
