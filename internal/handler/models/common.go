package models

type (
	ErrorResponse struct {
		Code    int64
		Reason  string
		Details string
	}
)
