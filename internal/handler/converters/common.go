package converters

import (
	"auth-service/internal/handler/models"
	"auth-service/internal/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeJsonErrorResponseWithErrorsError(err *errors.Error) models.ErrorResponse {
	res := models.ErrorResponse{
		Code:   err.Code,
		Reason: err.Reason,
	}

	if err.Details != nil {
		res.Details = err.Details.Error()
	}

	return res
}

func MakeStatusFromErrorsError(err *errors.Error) error {
	details := err.Reason

	if err.Details != nil {
		details = err.Details.Error()
	}

	return status.Errorf(codes.Internal, details)
}
