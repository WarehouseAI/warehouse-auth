package converters

import (
	"auth-service/internal/domain"
	"auth-service/internal/handler/models"
	"auth-service/internal/warehousepb"
)

func DomainCreateUser2ProtoCreateUser(request models.CreateRequestData) *warehousepb.CreateUserRequest {
	return &warehousepb.CreateUserRequest{
		Username:  request.Username,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Hash:      request.Password,
		Email:     request.Email,
	}
}

func DomainResetPassword2ProtoResetPassword(request domain.ResetPasswordRequestData) *warehousepb.ResetPasswordRequest {
	return &warehousepb.ResetPasswordRequest{
		UserId:   request.Id,
		Password: request.Password,
	}
}

func DomainUpdateVerification2ProtoUpdateVerification(request domain.UpdateVerificationRequestData) *warehousepb.UpdateVerificationStatusRequest {
	return &warehousepb.UpdateVerificationStatusRequest{
		UserId: request.Id,
		Email:  request.Email,
	}
}

func ProtoUser2DomainAccount(response *warehousepb.User) domain.Account {
	return domain.Account{
		Id:       response.UserId,
		Role:     domain.Role(response.Role),
		Verified: response.Verified,
	}
}

func DomainUser2ProtoAccount(acc domain.Account) *warehousepb.User {
	return &warehousepb.User{
		UserId:    acc.Id,
		Role:      int64(acc.Role),
		Username:  acc.Username,
		Firstname: acc.Firstname,
		Verified:  acc.Verified,
		Email:     acc.Email,
	}
}
