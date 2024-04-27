package jwt

import (
	randomAdpt "auth-service/internal/adapter/random"
	timeAdpt "auth-service/internal/adapter/time"
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"auth-service/internal/pkg/errors"
	"auth-service/internal/pkg/logger"
	jwtRepo "auth-service/internal/repository/operations/jwt"
	"auth-service/internal/repository/operations/transactions"
	"context"
	"sync"
	"time"
)

type (
	Service interface {
		Auth(ctx context.Context, token string, purpose domain.AuthPurpose) (domain.Account, int64, *errors.Error)
		Logout(ctx context.Context, role domain.Role, userId string) *errors.Error
		CreateTokens(ctx context.Context, role domain.Role, userId string) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error)
		CreateTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error)
		ReCreateTokens(ctx context.Context, role domain.Role, userId string, number int64) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error)
		DropTokens(ctx context.Context, role domain.Role, userId string, number int64) *errors.Error
		DropOldTokens(ctx context.Context, timestamp int64) *errors.Error
	}

	service struct {
		log logger.Logger

		txRepo transactions.Repository
		repo   jwtRepo.Repository

		timeAdapter   timeAdpt.Adapter
		randomAdapter randomAdpt.Adapter

		jwtKey                            string
		atTimeout, rtTimeout, authTimeout time.Duration
		lock                              *sync.RWMutex
	}
)

func NewService(
	log logger.Logger,
	txRepo transactions.Repository,
	repo jwtRepo.Repository,
	cfg config.Auth,
	timeAdapter timeAdpt.Adapter,
	randomAdapter randomAdpt.Adapter,
) Service {
	return &service{
		log:           log,
		txRepo:        txRepo,
		repo:          repo,
		lock:          &sync.RWMutex{},
		jwtKey:        cfg.Key,
		atTimeout:     cfg.AccessTokenTimeout,
		rtTimeout:     cfg.RefreshTokenTimeout,
		authTimeout:   cfg.AuthTimeout,
		timeAdapter:   timeAdapter,
		randomAdapter: randomAdapter,
	}
}

func (s *service) Auth(
	ctx context.Context, token string, purpose domain.AuthPurpose,
) (domain.Account, int64, *errors.Error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx, e := s.txRepo.StartTransaction(ctx)
	if e != nil {
		return domain.Account{}, 0, s.log.ServiceTxError(e)
	}
	defer tx.Rollback()

	t, err := s.parseToken(token)
	if err != nil {
		return domain.Account{}, 0, s.log.ServiceError(err)
	}

	acc, number, err := s.checkToken(ctx, tx, t, purpose)
	if err != nil {
		return domain.Account{}, 0, s.log.ServiceError(err)
	}

	if e = tx.Commit(); e != nil {
		return domain.Account{}, 0, s.log.ServiceTxError(e)
	}

	return acc, number, nil
}

func (s *service) Logout(ctx context.Context, role domain.Role, userId string) *errors.Error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx, e := s.txRepo.StartTransaction(ctx)
	if e != nil {
		return s.log.ServiceTxError(e)
	}
	defer tx.Rollback()

	if err := s.repo.DropAllTokensTX(ctx, tx, role, userId); err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	if e = tx.Commit(); e != nil {
		return s.log.ServiceTxError(e)
	}

	return nil
}

func (s *service) CreateTokens(ctx context.Context, role domain.Role, userId string) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	_, accessToken, refreshToken, e := s.createTokens(ctx, tx, role, userId)
	if e != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceError(e)
	}

	if err = tx.Commit(); err != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) CreateTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, accessToken, refreshToken, e := s.createTokens(ctx, tx, role, userId)
	if e != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceError(e)
	}

	return accessToken, refreshToken, nil
}

func (s *service) ReCreateTokens(ctx context.Context, role domain.Role, userId string, number int64) (domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	if err := s.repo.DropTokensTX(ctx, tx, role, userId, number); err != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}

	_, accessToken, refreshToken, e := s.createTokens(ctx, tx, role, userId)
	if e != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceError(errors.WD(errors.AuthCreateTokens, err))
	}

	if err = tx.Commit(); err != nil {
		return domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) DropTokens(ctx context.Context, role domain.Role, userId string, number int64) *errors.Error {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	err = s.repo.DropTokensTX(ctx, tx, role, userId, number)
	if err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	if err = tx.Commit(); err != nil {
		return s.log.ServiceTxError(err)
	}

	return nil
}

func (s *service) DropOldTokens(ctx context.Context, timestamp int64) *errors.Error {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	err = s.repo.DropOldTokens(ctx, tx, timestamp)
	if err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	if err = tx.Commit(); err != nil {
		return s.log.ServiceTxError(err)
	}

	return nil
}
