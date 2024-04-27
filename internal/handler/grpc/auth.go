package grpc

import (
	"auth-service/internal/converters"
	"auth-service/internal/domain"
	handler_converters "auth-service/internal/handler/converters"
	"auth-service/internal/pkg/logger"
	"auth-service/internal/service/jwt"
	"auth-service/internal/warehousepb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	warehousepb.UnimplementedAuthServer
	log    logger.Logger
	jwtSvc jwt.Service
}

func (s *GrpcServer) Authenticate(ctx context.Context, req *warehousepb.AuthRequest) (*warehousepb.AuthResponse, error) {
	if req == nil || req.AccessToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	acc, _, err := s.jwtSvc.Auth(ctx, req.AccessToken, domain.PurposeAccess)

	if err != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(err)
	}

	return &warehousepb.AuthResponse{User: converters.DomainUser2ProtoAccount(acc)}, nil
}
