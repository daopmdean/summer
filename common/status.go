package common

type ResponseStatusValue string

type responseStatus struct {
	Success      ResponseStatusValue
	NotFound     ResponseStatusValue
	Invalid      ResponseStatusValue
	Unauthorized ResponseStatusValue
	Error        ResponseStatusValue
}

var ResponseStatus = &responseStatus{
	Success:      "SUCCESS",
	NotFound:     "NOT_FOUND",
	Invalid:      "INVALID",
	Unauthorized: "UNAUTHORIZED",
	Error:        "ERROR",
}
