package auth

import (
	mailAdpt "auth-service/internal/adapter/mail"
	timeAdpt "auth-service/internal/adapter/time"
	userAdpt "auth-service/internal/adapter/user"
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"auth-service/internal/handler/models"
	"auth-service/internal/pkg/errors"
	"auth-service/internal/pkg/logger"
	"auth-service/internal/pkg/utils/encode"
	"auth-service/internal/pkg/utils/str"
	rep_converters "auth-service/internal/repository/converters"
	jwtRepo "auth-service/internal/repository/operations/jwt"
	"auth-service/internal/repository/operations/reset_token"
	"auth-service/internal/repository/operations/transactions"
	"auth-service/internal/repository/operations/verification_token"
	jwtSvc "auth-service/internal/service/jwt"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO:
// - Add confirm verification token
// - Add confirm reset token
type (
	Service interface {
		Login(ctx context.Context, reqData models.LoginRequestData) (*domain.Account, domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error)
		Register(ctx context.Context, reqData models.CreateRequestData) (string, *errors.Error)
		FullLogout(ctx context.Context, role domain.Role, accId string) *errors.Error
		CheckVerificationToken(ctx context.Context, vt, accId, tokenId string) (domain.Account, *errors.Error)
		CreateResetToken(ctx context.Context, email string) *errors.Error
		VerifyResetToken(ctx context.Context, token, tokenId, accId string) *errors.Error
	}

	service struct {
		cfg config.Config

		txRepo           transactions.Repository
		jwtRepo          jwtRepo.Repository
		verificationRepo verification_token.Repository
		resetRepo        reset_token.Repository

		log        logger.Logger
		jwtService jwtSvc.Service

		timeAdapter timeAdpt.Adapter
		userAdapter userAdpt.Adapter
		mailAdapter mailAdpt.Adapter
	}
)

func NewService(
	cfg config.Config,
	txRepo transactions.Repository,
	jwt jwtRepo.Repository,
	timeAdapter timeAdpt.Adapter,
	jwtService jwtSvc.Service,
	log logger.Logger,
	userAdapter userAdpt.Adapter,
	mailAdapter mailAdpt.Adapter,
	verificationRepo verification_token.Repository,
	resetRepo reset_token.Repository,
) Service {
	return &service{
		cfg:              cfg,
		txRepo:           txRepo,
		jwtRepo:          jwt,
		timeAdapter:      timeAdapter,
		jwtService:       jwtService,
		log:              log,
		userAdapter:      userAdapter,
		mailAdapter:      mailAdapter,
		verificationRepo: verificationRepo,
		resetRepo:        resetRepo,
	}
}

func (s *service) VerifyResetToken(ctx context.Context, token, tokenId, accId string) *errors.Error {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	// Check user is exists
	if _, err = s.userAdapter.GetById(ctx, accId); err != nil {
		return s.log.ServiceGrpcAdapterError(err)
	}

	hashedToken, err := encode.HashedPassword(token)
	if err != nil {
		return s.log.ServiceError(errors.WD(errors.InternalError, err))
	}

	rt, err := s.resetRepo.GetById(ctx, tx, tokenId)
	if err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	if rt.ExpiresAt < time.Now().Unix() {
		return errors.AuthExpiredToken
	}

	if rt.Token != hashedToken {
		return errors.AuthInvalidToken
	}

	if err := s.resetRepo.DeleteById(ctx, tx, tokenId); err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return s.log.ServiceTxError(err)
	}

	return nil
}

func (s *service) CreateResetToken(ctx context.Context, email string) *errors.Error {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	acc, err := s.userAdapter.GetByEmail(ctx, email)
	if err != nil {
		return s.log.ServiceGrpcAdapterError(err)
	}

	token := str.RandomString(16)
	hashedToken, err := encode.HashedPassword(token)
	if err != nil {
		return s.log.ServiceError(errors.WD(errors.InternalError, err))
	}

	rtInfo := domain.ResetTokenInfo{
		UserId:    acc.Id,
		Token:     hashedToken,
		CreatedAt: s.timeAdapter.Now().Unix(),
		ExpiresAt: s.timeAdapter.AddTime(s.timeAdapter.Now(), time.Minute*15).Unix(),
	}

	rt, err := s.resetRepo.Create(ctx, tx, rep_converters.DomainResetToken2ModelResetToken(rtInfo))
	if err != nil {
		return s.log.ServiceDatabaseError(err)
	}

	mail := domain.EmailMessage{
		To:   acc.Email,
		Type: domain.ResetType,
		Payload: domain.Payload{
			Firstname: acc.Firstname,
			ResetPayload: domain.ResetPayload{
				TokenId: rt.ID.String(),
				Token:   token,
				AccId:   acc.Id,
			},
		},
	}

	if err := s.mailAdapter.SendMessage(mail); err != nil {
		return s.log.ServiceBrokerAdapterError(err)
	}

	if err := tx.Commit(); err != nil {
		return s.log.ServiceTxError(err)
	}

	return nil
}

func (s *service) FullLogout(ctx context.Context, role domain.Role, userId string) *errors.Error {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	if err := s.jwtRepo.DropAllTokensTX(ctx, tx, role, userId); err != nil {
		return errors.DatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return s.log.ServiceTxError(err)
	}

	return nil
}

func (s *service) CheckVerificationToken(ctx context.Context, vt string, accId string, tokenId string) (domain.Account, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.Account{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	hashedToken, err := encode.HashedPassword(vt)
	if err != nil {
		return domain.Account{}, s.log.ServiceError(errors.WD(errors.InternalError, err))
	}

	acc, err := s.userAdapter.GetById(ctx, accId)
	if err != nil {
		return domain.Account{}, s.log.ServiceGrpcAdapterError(err)
	}
	info, err := s.verificationRepo.GetById(ctx, tx, tokenId)
	if err != nil {
		return domain.Account{}, s.log.ServiceDatabaseError(err)
	}
	if info.Token != hashedToken {
		return domain.Account{}, errors.AuthInvalidToken
	}

	success, err := s.userAdapter.UpdateVerificationStatus(ctx, domain.UpdateVerificationRequestData{Id: accId, Email: acc.Email})
	if err != nil {
		return domain.Account{}, s.log.ServiceGrpcAdapterError(err)
	}
	if !success {
		return domain.Account{}, errors.AuthVerificationFailed
	}

	acc.Verified = true
	if err := tx.Commit(); err != nil {
		return domain.Account{}, s.log.ServiceTxError(err)
	}

	return acc, nil
}

func (s *service) Login(
	ctx context.Context, reqData models.LoginRequestData,
) (*domain.Account, domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return nil, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	acc, err := s.userAdapter.GetByLogin(ctx, reqData.Login)
	if err != nil {
		return nil, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceGrpcAdapterError(err)
	}

	if !acc.Verified {
		return nil, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, errors.AuthNotVerifiedAccount
	}

	accessToken, refreshToken, e := s.jwtService.CreateTokensTX(ctx, tx, acc.Role, acc.Id)
	if e != nil {
		return nil, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, e
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}

	return &acc, accessToken, refreshToken, nil
}

func (s *service) Register(
	ctx context.Context, reqData models.CreateRequestData,
) (string, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return "", s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	_, err = s.userAdapter.GetByEmail(ctx, reqData.Email)
	if err != nil {
		extError := status.Convert(err)
		if extError.Code() != codes.AlreadyExists {
			return "", s.log.ServiceGrpcAdapterError(err)
		}

		return "", errors.AuthUserAlreadyExists
	}

	_, err = s.userAdapter.GetByLogin(ctx, reqData.Username)
	if err != nil {
		extError := status.Convert(err)
		if extError.Code() != codes.AlreadyExists {
			return "", s.log.ServiceGrpcAdapterError(err)
		}

		return "", errors.AuthUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), 12)
	if err != nil {
		return "", s.log.ServiceError(errors.WD(errors.InternalError, err))
	}
	reqData.Password = string(hash)

	acc, err := s.userAdapter.CreateUser(ctx, reqData)
	if err != nil {
		return "", s.log.ServiceGrpcAdapterError(err)
	}

	token := str.RandomString(6)
	hashedToken, err := encode.HashedPassword(token)
	if err != nil {
		return "", s.log.ServiceError(errors.WD(errors.InternalError, err))
	}

	vtInfo := domain.VerificationTokenInfo{
		UserId:    acc.Id,
		Token:     hashedToken,
		SendTo:    acc.Email,
		CreatedAt: s.timeAdapter.Now().Unix(),
		ExpiresAt: s.timeAdapter.AddTime(s.timeAdapter.Now(), time.Minute*10).Unix(),
	}
	vt, err := s.verificationRepo.Create(
		ctx, tx,
		rep_converters.DomainVerificationToken2ModelVerificationToken(vtInfo),
	)
	if err != nil {
		return "", s.log.ServiceDatabaseError(err)
	}

	email := domain.EmailMessage{
		Type: domain.VerificationType,
		To:   acc.Email,
		Payload: domain.Payload{
			Firstname: reqData.Firstname,
			VerifyPayload: domain.VerifyPayload{
				Token: token,
			},
		},
	}

	if err := s.mailAdapter.SendMessage(email); err != nil {
		return "", s.log.ServiceError(errors.WD(errors.InternalError, err))
	}

	return vt.ID.String(), nil
}
