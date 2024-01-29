package server

import (
	"auth-service/internal/app/dataservice"
	"auth-service/internal/app/service"
	e "auth-service/internal/pkg/errors/grpc"
	gen "auth-service/internal/pkg/protogen"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	gen.UnimplementedAuthServiceServer
	DB     dataservice.SessionInterface
	Logger *logrus.Logger
}

func (s *AuthGrpcServer) Authenticate(ctx context.Context, req *gen.AuthenticationRequest) (*gen.AuthenticationResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	userId, session, err := service.Authenticate(req.SessionId, s.DB)

	if err != nil {
		return nil, e.NewGrpcErrorByHttpStatus(err)
	}

	return &gen.AuthenticationResponse{UserId: *userId, SessionId: session.ID}, nil
}
