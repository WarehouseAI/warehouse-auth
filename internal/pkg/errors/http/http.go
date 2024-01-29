package errors

import (
	dbe "auth-service/internal/pkg/errors/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatusCode int

type Error interface {
	error
	GetStatus() int
	GetTrace() string
}

type HttpError struct {
	Status int
	Trace  string
	Err    error
}

func (he HttpError) Error() string {
	return he.Err.Error()
}

func (he HttpError) GetStatus() int {
	return he.Status
}

func (he HttpError) GetTrace() string {
	return he.Trace
}

func NewHttpErrorByDbStatus(dbErr error) error {
	err := dbErr.(dbe.DBError)

	switch err.Status {
	case dbe.DbExist:
		return HttpError{409, err.Trace, dbErr}

	case dbe.DbNotFound:
		return HttpError{404, err.Trace, dbErr}

	default:
		return HttpError{500, err.Trace, dbErr}
	}
}

func NewHttpErrorByGrpcStatus(grpcErr error) error {
	err, _ := status.FromError(grpcErr)

	switch err.Code() {
	case codes.InvalidArgument:
		return HttpError{400, err.Message(), grpcErr}

	case codes.Unauthenticated:
		return HttpError{401, err.Message(), grpcErr}

	case codes.PermissionDenied:
		return HttpError{403, err.Message(), grpcErr}

	case codes.NotFound:
		return HttpError{404, err.Message(), grpcErr}

	case codes.Aborted:
		return HttpError{409, err.Message(), grpcErr}

	default:
		return HttpError{500, err.Message(), grpcErr}
	}
}

func NewHttpError(statusCode int, trace string, err error) error {
	return HttpError{statusCode, trace, err}
}
