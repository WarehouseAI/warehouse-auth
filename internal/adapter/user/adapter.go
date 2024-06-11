package user

import (
	"context"

	"github.com/warehouse/auth-service/internal/config"
	"github.com/warehouse/auth-service/internal/converters"
	"github.com/warehouse/auth-service/internal/domain"
	"github.com/warehouse/auth-service/internal/handler/models"
	"github.com/warehouse/auth-service/internal/warehousepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	Adapter interface {
		CreateUser(ctx context.Context, request models.CreateRequestData) (domain.Account, error)
		ResetPassword(ctx context.Context, request domain.ResetPasswordRequestData) (bool, error)
		GetByEmail(ctx context.Context, email string) (domain.Account, error)
		GetByLogin(ctx context.Context, username string) (domain.Account, error)
		GetById(ctx context.Context, userId string) (domain.Account, error)
		UpdateVerificationStatus(ctx context.Context, request domain.UpdateVerificationRequestData) (bool, error)
	}

	adapter struct {
		client warehousepb.UserServiceClient
		config config.Grpc
	}
)

func NewAdapter(
	config config.Grpc,
) (Adapter, error) {
	conn, err := grpc.NewClient(config.User.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := warehousepb.NewUserServiceClient(conn)

	return &adapter{
		client: client,
		config: config,
	}, nil
}

func (a *adapter) CreateUser(ctx context.Context, request models.CreateRequestData) (domain.Account, error) {
	resp, err := a.client.CreateUser(ctx, converters.DomainCreateUser2ProtoCreateUser(request))

	if err != nil {
		return domain.Account{}, err
	}

	return converters.ProtoUser2DomainAccount(resp), nil
}

func (a *adapter) ResetPassword(ctx context.Context, request domain.ResetPasswordRequestData) (bool, error) {
	resp, err := a.client.ResetPassword(ctx, converters.DomainResetPassword2ProtoResetPassword(request))

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (a *adapter) UpdateVerificationStatus(ctx context.Context, request domain.UpdateVerificationRequestData) (bool, error) {
	resp, err := a.client.UpdateVerificationStatus(ctx, converters.DomainUpdateVerification2ProtoUpdateVerification(request))

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (a *adapter) GetById(ctx context.Context, userId string) (domain.Account, error) {
	resp, err := a.client.GetUserById(ctx, &warehousepb.GetUserByIdRequest{Id: userId})

	if err != nil {
		return domain.Account{}, err
	}

	return converters.ProtoUser2DomainAccount(resp), nil
}

func (a *adapter) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	resp, err := a.client.GetUserByEmail(ctx, &warehousepb.GetUserByEmailRequest{Email: email})

	if err != nil {
		return domain.Account{}, err
	}

	return converters.ProtoUser2DomainAccount(resp), nil
}

func (a *adapter) GetByLogin(ctx context.Context, username string) (domain.Account, error) {
	resp, err := a.client.GetUserByLogin(ctx, &warehousepb.GetUserByLoginRequest{Username: username})

	if err != nil {
		return domain.Account{}, err
	}

	return converters.ProtoUser2DomainAccount(resp), nil
}
