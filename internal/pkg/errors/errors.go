package errors

import (
	"errors"

	wh_converters "github.com/warehouse/auth-service/internal/pkg/utils/converters"
)

type Error struct {
	Code    int64  `json:"code"`
	Reason  string `json:"reason"`
	Details error  `json:"details"`
}

var (
	HashPasswordRaw = errors.New("hash password error")
	HashPassword    = &Error{Code: 500, Reason: "hash password error"}

	ServiceUnavailable = &Error{Code: 500, Reason: "service unavailable"}
	HttpRequestError   = &Error{Reason: "http request error"}

	PermissionDenied   = &Error{Code: 403, Reason: "permission denied"}
	ParseError         = &Error{Code: 400, Reason: "parse failed"}
	InternalError      = &Error{Code: 500, Reason: "internal error"}
	MissingCredentials = &Error{Code: 401, Reason: "missing credentials"}
	ValidationFailed   = &Error{Code: 400, Reason: "validation failed"}
	Timeout            = &Error{Code: 504, Reason: "timeout"}
	InvalidVersion     = &Error{Code: 400, Reason: "invalid version"}
)

func New(text string) error {
	return errors.New(text)
}

// WD WD значит WithDetails
func WD(err *Error, details error) *Error {
	e := *err
	e.Details = details
	return &e
}

func DatabaseError(details error) *Error {
	return WD(&Error{Code: 500, Reason: "database failed"}, details)
}

func BrokerError(details error) *Error {
	return WD(&Error{Code: 500, Reason: "broker error"}, details)
}

func GrpcError(details error) *Error {
	return WD(&Error{Code: int64(wh_converters.GrpcStatusToHttp(details)), Reason: "grpc request failed"}, details)
}

func WithDetailsAndCode(err *Error, details error, code int64) *Error {
	e := *err
	e.Details = details
	e.Code = code
	return &e
}

func Equal(e1, e2 *Error) bool {
	return (e1 == nil && e2 == nil) ||
		(e1 != nil && e2 != nil && e1.Code == e2.Code && e1.Reason == e2.Reason && e1.Details == e2.Details)
}
