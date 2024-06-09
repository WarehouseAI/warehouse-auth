package http

import (
	"context"
	"encoding/json"
	"net/http"

	timeAdpt "github.com/warehouse/auth-service/internal/adapter/time"
	userAdpt "github.com/warehouse/auth-service/internal/adapter/user"
	"github.com/warehouse/auth-service/internal/config"
	"github.com/warehouse/auth-service/internal/domain"
	"github.com/warehouse/auth-service/internal/handler/middlewares"
	"github.com/warehouse/auth-service/internal/handler/models"
	"github.com/warehouse/auth-service/internal/pkg/errors"
	"github.com/warehouse/auth-service/internal/service/auth"
	"github.com/warehouse/auth-service/internal/service/jwt"

	"github.com/gorilla/mux"
)

type (
	authHandler struct {
		cfg      *config.Server
		timeouts *config.Timeouts

		jwtService  jwt.Service
		authService auth.Service

		timeAdapter timeAdpt.Adapter
		userAdapter userAdpt.Adapter

		reqHandler WarehouseRequestHandler
		middleware middlewares.Middleware
	}
)

func NewAuthHandler(
	cfg config.Server,
	timeouts config.Timeouts,

	jwtSvc jwt.Service,
	authSvc auth.Service,

	timeAdpt timeAdpt.Adapter,
	userAdpt userAdpt.Adapter,

	requestHandler WarehouseRequestHandler,
	middlewares middlewares.Middleware,
) Handler {
	return &authHandler{
		cfg:      &cfg,
		timeouts: &timeouts,

		jwtService:  jwtSvc,
		authService: authSvc,

		timeAdapter: timeAdpt,
		userAdapter: userAdpt,

		reqHandler: requestHandler,
		middleware: middlewares,
	}
}

func (h *authHandler) Shutdown() {
}

func (h *authHandler) FillHandlers(router *mux.Router) {
	base := "/auth"
	r := router.PathPrefix(base).Subrouter()
	h.reqHandler.HandleJsonRequest(r, base, "", http.MethodPost, h.loginHandler)
	h.reqHandler.HandleJsonRequestWithMiddleware(r, base, "/full_logout", http.MethodDelete, h.fullLogoutHandler, h.middleware.JwtAuthMiddleware(domain.PurposeAccess))
	h.reqHandler.HandleJsonRequestWithMiddleware(r, base, "/logout", http.MethodDelete, h.logoutHandler, h.middleware.JwtAuthMiddleware(domain.PurposeAccess))
	h.reqHandler.HandleJsonRequestWithMiddleware(r, base, "/refresh", http.MethodGet, h.refreshHandler, h.middleware.JwtAuthMiddleware(domain.PurposeRefresh))
	h.reqHandler.HandleJsonRequest(r, base, "/register", http.MethodPost, h.registerHandler)
	h.reqHandler.HandleJsonRequest(r, base, "/verify/check", http.MethodGet, h.checkVerificationToken)
	h.reqHandler.HandleJsonRequest(r, base, "/reset/request", http.MethodGet, h.resetPasswordRequest)
	h.reqHandler.HandleJsonRequest(r, base, "/reset/confirm", http.MethodPost, h.resetPasswordConfirm)
}

func (h *authHandler) refreshHandler(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	if acc == nil {
		return whJsonErrorResponse(errors.AuthAuthFailed)
	}
	number := ctx.Value(domain.TokenNumberCtxKey).(int64)

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()
	newAt, newRt, err := h.jwtService.ReCreateTokens(ctx, acc.Role, acc.Id, number)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	accCookie, err := createCookie("acc", acc)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		models.Tokens{
			AccessToken:  newAt,
			RefreshToken: newRt,
		},
		http.StatusOK,
		[]http.Cookie{accCookie},
	)
}

func (h *authHandler) fullLogoutHandler(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	if acc == nil {
		return whJsonErrorResponse(errors.AuthAuthFailed)
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	if err := h.authService.FullLogout(ctx, acc.Role, acc.Id); err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		nil,
		http.StatusNoContent,
		nil,
	)
}

func (h *authHandler) logoutHandler(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	if acc == nil {
		return whJsonErrorResponse(errors.AuthAuthFailed)
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	if err := h.jwtService.Logout(ctx, acc.Role, acc.Id); err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		nil,
		http.StatusNoContent,
		nil,
	)
}

func (h *authHandler) loginHandler(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	ctx, cancel := context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	var req models.LoginRequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return whJsonErrorResponse(errors.WD(errors.InternalError, err))
	}

	acc, at, rt, err := h.authService.Login(ctx, req)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	accCookie, err := createCookie("acc", acc)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		models.Tokens{
			AccessToken:  at,
			RefreshToken: rt,
		},
		http.StatusOK,
		[]http.Cookie{accCookie},
	)
}

func (h *authHandler) registerHandler(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	ctx, cancel := context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	var req models.CreateRequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return whJsonErrorResponse(errors.WD(errors.InternalError, err))
	}

	verificationTokenId, err := h.authService.Register(ctx, req)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		models.CreateResponsedata{
			VerificationTokenId: verificationTokenId,
		},
		http.StatusCreated,
		nil,
	)
}

func (h *authHandler) checkVerificationToken(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	ctx, cancel := context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	vars := mux.Vars(r)
	plainVerificationToken := vars["token"]
	accId := vars["acc_id"]
	tokenId := vars["token_id"]

	existAcc, err := h.authService.CheckVerificationToken(ctx, plainVerificationToken, accId, tokenId)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	accessToken, refreshToken, err := h.jwtService.CreateTokens(ctx, existAcc.Role, accId)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	accCookie, err := createCookie("acc", existAcc)
	if err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		models.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		http.StatusOK,
		[]http.Cookie{accCookie},
	)
}

func (h *authHandler) resetPasswordRequest(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	ctx, cancel := context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	accEmail := mux.Vars(r)["email"]

	if err := h.authService.CreateResetToken(ctx, accEmail); err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		nil,
		http.StatusOK,
		nil,
	)
}

func (h *authHandler) resetPasswordConfirm(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse {
	ctx, cancel := context.WithTimeout(ctx, h.timeouts.RequestTimeout)
	defer cancel()

	var req models.PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return whJsonErrorResponse(errors.WD(errors.InternalError, err))
	}

	if err := h.authService.VerifyResetToken(ctx, req.Token, req.TokenId, req.AccId); err != nil {
		return whJsonErrorResponse(err)
	}

	return whJsonSuccessResponse(
		nil,
		http.StatusOK,
		nil,
	)
}
