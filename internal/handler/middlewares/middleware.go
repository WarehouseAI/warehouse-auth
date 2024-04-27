package middlewares

import (
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"auth-service/internal/pkg/logger"
	"auth-service/internal/service/jwt"
	"net/http"
)

const (
	VersionDelimiter = ":" // Разделитель составных частей версий
	VersionHeader    = "Coffee-Version"

	AuthHeader    = "Authorization"
	TokenStart    = "Bearer "       // Префикс значения заголовка с авторизацией
	TokenStartInd = len(TokenStart) // Индекс, с которого в заголовке авторизации должен начинаться jwt токен

	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

type (
	Middleware interface {
		JwtAuthMiddleware(purpose domain.AuthPurpose) func(http.Handler) http.Handler
		QueueMiddleware(h http.Handler) http.Handler
	}

	middleware struct {
		log logger.Logger

		timeouts   config.Timeouts
		jwtService jwt.Service
		queue      chan struct{}
	}
)

func NewMiddleware(
	log logger.Logger,
	timeouts config.Timeouts,
	jwtService jwt.Service,
) Middleware {
	return &middleware{
		log:        log,
		timeouts:   timeouts,
		jwtService: jwtService,
	}
}
