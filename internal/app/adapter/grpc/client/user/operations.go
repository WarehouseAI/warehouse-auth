package user

import (
	e "auth-service/internal/pkg/errors/http"
	gen "auth-service/internal/pkg/protogen"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcClient struct {
	conn gen.UserServiceClient
}

func NewUserGrpcClient(grpcUrl string) *UserGrpcClient {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	client := gen.NewUserServiceClient(conn)

	return &UserGrpcClient{
		conn: client,
	}
}

func (c *UserGrpcClient) Create(ctx context.Context, userInfo *gen.CreateUserMsg) (string, error) {
	resp, err := c.conn.CreateUser(ctx, userInfo)

	if err != nil {
		return "", e.NewHttpErrorByGrpcStatus(err)
	}

	return resp.UserId, nil
}

func (c *UserGrpcClient) ResetPassword(ctx context.Context, resetPasswordRequest *gen.ResetPasswordRequest) (string, error) {
	resp, err := c.conn.ResetPassword(ctx, resetPasswordRequest)

	if err != nil {
		return "", e.NewHttpErrorByGrpcStatus(err)
	}

	return resp.UserId, nil
}

func (c *UserGrpcClient) GetByEmail(ctx context.Context, email string) (*gen.User, error) {
	resp, err := c.conn.GetUserByEmail(ctx, &gen.GetUserByEmailMsg{Email: email})

	if err != nil {
		return nil, e.NewHttpErrorByGrpcStatus(err)
	}

	return resp, nil
}

func (c *UserGrpcClient) UpdateVerificationStatus(ctx context.Context, userId string) (bool, error) {
	resp, err := c.conn.UpdateVerificationStatus(ctx, &gen.UpdateVerificationStatusRequest{UserId: userId})

	if err != nil {
		return false, e.NewHttpErrorByGrpcStatus(err)
	}

	return resp.Verified, nil
}

func (c *UserGrpcClient) GetById(ctx context.Context, userId string) (*gen.User, error) {
	resp, err := c.conn.GetUserById(ctx, &gen.GetUserByIdMsg{UserId: userId})

	if err != nil {
		return nil, e.NewHttpErrorByGrpcStatus(err)
	}

	return resp, nil
}
