package grpc

import (
	e "auth-service/internal/pkg/errors/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpcErrorByHttpStatus(httpErr error) error {
	err := httpErr.(e.HttpError)

	switch err.Status {
	// Bad Request, Unprocessable Entity
	case 400, 422:
		return status.Errorf(codes.InvalidArgument, httpErr.Error())

	// Unauthorized
	case 401:
		return status.Errorf(codes.Unauthenticated, httpErr.Error())

	// Forbidden
	case 403:
		return status.Errorf(codes.PermissionDenied, httpErr.Error())

	// Not Found
	case 404:
		return status.Errorf(codes.NotFound, httpErr.Error())

	// Conflict
	case 409:
		return status.Errorf(codes.Aborted, httpErr.Error())

	default:
		return status.Errorf(codes.Internal, httpErr.Error())
	}
}
