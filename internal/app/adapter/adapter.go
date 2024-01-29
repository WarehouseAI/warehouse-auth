package adapter

import (
	gen "auth-service/internal/pkg/protogen"
	"context"
)

type UserGrpcInterface interface {
	Create(ctx context.Context, userInfo *gen.CreateUserMsg) (string, error)
	ResetPassword(ctx context.Context, resetPasswordRequest *gen.ResetPasswordRequest) (string, error)
	GetByEmail(ctx context.Context, email string) (*gen.User, error)
	GetById(ctx context.Context, userId string) (*gen.User, error)
	UpdateVerificationStatus(ctx context.Context, userId string) (bool, error)
}

type BrokerInterface interface {
	SendMessage(queueName string, message interface{}) error
}
