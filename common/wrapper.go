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
	return BuildRes(ResponseStatus.Invalid, errCode, errMsg)
}

func BuildErrorRes(errCode, errMsg string) *Response {
	return BuildRes(ResponseStatus.Error, errCode, errMsg)
}

func BuildNotfoundRes(errCode, errMsg string) *Response {
	return BuildRes(ResponseStatus.NotFound, errCode, errMsg)
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

func BuildSingleListError(errCode, errMsg string) []*ErrRes {
	return []*ErrRes{
		BuildSingleError(errCode, errMsg),
	}
}

func BuildSingleError(errCode, errMsg string) *ErrRes {
	return &ErrRes{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}
