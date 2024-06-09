package grpc

import (
	"context"

	"github.com/warehouse/auth-service/internal/config"
	"github.com/warehouse/auth-service/internal/converters"
	"github.com/warehouse/auth-service/internal/domain"
	handler_converters "github.com/warehouse/auth-service/internal/handler/converters"
	"github.com/warehouse/auth-service/internal/pkg/logger"
	"github.com/warehouse/auth-service/internal/service/jwt"
	"github.com/warehouse/auth-service/internal/warehousepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	AuthHandler struct {
		warehousepb.UnimplementedAuthServer
		timeouts config.Timeouts
		log      logger.Logger
		jwtSvc   jwt.Service
	}
)

func NewAuthHandler(
	timeouts config.Timeouts,
	log logger.Logger,
	jwtSvc jwt.Service,
) *AuthHandler {
	return &AuthHandler{
		timeouts: timeouts,
		log:      log,
		jwtSvc:   jwtSvc,
	}
}

func (s *AuthHandler) Authenticate(ctx context.Context, req *warehousepb.AuthRequest) (*warehousepb.AuthResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.timeouts.RequestTimeout)
	defer cancel()

	acc, num, err := s.jwtSvc.Auth(ctx, req.Token, domain.AuthPurpose(req.Purpose))

	if err != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(err)
	}

	return &warehousepb.AuthResponse{User: converters.DomainUser2ProtoAccount(acc), Number: num}, nil
}
