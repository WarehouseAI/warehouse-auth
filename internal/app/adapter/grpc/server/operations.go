package server

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	"auth-service/internal/app/service"
	"auth-service/internal/app/service/confirm"
	"auth-service/internal/pkg/email"
	e "auth-service/internal/pkg/errors/grpc"
	gen "auth-service/internal/pkg/protogen"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	gen.UnimplementedAuthServiceServer
	SessionRepo dataservice.SessionInterface
	TokenRepo   dataservice.VerificationTokenInterface
	Broker      adapter.BrokerInterface
}

func (s *AuthGrpcServer) Authenticate(ctx context.Context, req *gen.AuthenticationRequest) (*gen.AuthenticationResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	userId, session, err := service.Authenticate(req.SessionId, s.SessionRepo)

	if err != nil {
		return nil, e.NewGrpcErrorByHttpStatus(err)
	}

	return &gen.AuthenticationResponse{UserId: *userId, SessionId: session.ID}, nil
}

func (s *AuthGrpcServer) SendEmailVerification(ctx context.Context, req *gen.VerifyEmailRequest) (*gen.VerifyEmailResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	verificationToken, err := confirm.CreateAndStoreVerificationToken(req.UserId, s.TokenRepo)

	if err != nil {
		return nil, e.NewGrpcErrorByHttpStatus(err)
	}

	if err := email.SendVerification(req.UserId, verificationToken, req.UserFirstname, req.UserEmail, s.Broker); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.VerifyEmailResponse{UserId: *&req.UserId}, nil
}
