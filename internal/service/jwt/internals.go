package jwt

import (
	"auth-service/internal/domain"
	"auth-service/internal/pkg/errors"
	"auth-service/internal/pkg/errors/service_errors"
	"auth-service/internal/repository/models"
	"auth-service/internal/repository/operations/transactions"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (s *service) parseToken(token string) (*jwt.Token, *errors.Error) {
	res, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtKey), nil
	})
	if err != nil && res == nil {
		return nil, s.log.ServiceError(errors.WD(errors.AuthParseToken, err))
	}
	return res, nil
}

func (s *service) checkToken(
	ctx context.Context, tx transactions.Transaction, token *jwt.Token, purpose domain.AuthPurpose,
) (domain.Account, int64, *errors.Error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.Account{}, 0, errors.AuthParseToken
	}

	if !claims.VerifyExpiresAt(s.timeAdapter.Now().Unix(), true) {
		return domain.Account{}, 0, errors.AuthExpiredToken
	}

	if realPurpose, ok := claims["purpose"].(float64); !ok {
		return domain.Account{}, 0, errors.AuthInvalidToken
	} else if domain.AuthPurpose(realPurpose) != purpose {
		return domain.Account{}, 0, errors.AuthInvalidTokenPurpose
	}

	if !token.Valid {
		return domain.Account{}, 0, errors.AuthInvalidToken
	}

	user_id, err := s.parseTokenStringClaim(claims, "user_id")
	if err != nil {
		return domain.Account{}, 0, err
	}
	number, err := s.parseTokenIntClaim(claims, "number")
	if err != nil {
		return domain.Account{}, 0, err
	}
	role, err := s.parseTokenIntClaim(claims, "role")
	if err != nil {
		return domain.Account{}, 0, err
	}
	secret, err := s.parseTokenStringClaim(claims, "secret")
	if err != nil {
		return domain.Account{}, 0, err
	}

	if _, ok := s.repo.GetTokenMap()[domain.Role(role)]; !ok {
		return domain.Account{}, 0, errors.AuthInvalidToken
	}

	tokenModel := models.Token{
		UserId:  user_id,
		Number:  number,
		Purpose: int(purpose),
		Secret:  secret,
	}
	if _, e := s.repo.CheckTokenTX(ctx, tx, domain.Role(role), tokenModel); e != nil {
		return domain.Account{}, 0, errors.WD(service_errors.DatabaseError, e)
	}

	// TODO: добавить подтяг данных пользователя 
	return domain.Account{
		Role: domain.Role(role),
		Id:   user_id,
	}, number, nil
}

func (s *service) parseTokenIntClaim(claims jwt.MapClaims, key string) (int64, *errors.Error) {
	if parsedValue, ok := claims[key].(float64); !ok {
		return 0, errors.AuthInvalidToken
	} else {
		return int64(parsedValue), nil
	}
}

func (s *service) parseTokenStringClaim(claims jwt.MapClaims, key string) (string, *errors.Error) {
	if stringValue, ok := claims[key].(string); !ok {
		return "", errors.AuthInvalidToken
	} else {
		return stringValue, nil
	}
}

func (s *service) createTokens(
	ctx context.Context, tx transactions.Transaction, role domain.Role, userId string,
) (int64, domain.JwtTokenInfo, domain.JwtTokenInfo, *errors.Error) {
	now := s.timeAdapter.Now()

	number, err := s.repo.FindNumberTX(ctx, tx, role, userId)
	if err != nil {
		return 0, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, errors.DatabaseError(err)
	}
	accessExpiresAt, refreshExpiresAt := now.Add(s.atTimeout), now.Add(s.rtTimeout)

	accessTokenHash, err := s.generateTokenHash(ctx, tx, role, userId, number, domain.PurposeAccess, accessExpiresAt)
	if err != nil {
		return 0, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, errors.DatabaseError(err)
	}

	refreshTokenHash, err := s.generateTokenHash(ctx, tx, role, userId, number, domain.PurposeRefresh, refreshExpiresAt)
	if err != nil {
		return 0, domain.JwtTokenInfo{}, domain.JwtTokenInfo{}, s.log.ServiceTxError(err)
	}

	accessToken := domain.JwtTokenInfo{
		Token:     accessTokenHash,
		ExpiresAt: accessExpiresAt.UnixNano() / 1e+6,
	}
	refreshToken := domain.JwtTokenInfo{
		Token:     refreshTokenHash,
		ExpiresAt: refreshExpiresAt.UnixNano() / 1e+6,
	}

	return number, accessToken, refreshToken, nil
}

func (s *service) generateTokenHash(
	ctx context.Context, tx transactions.Transaction, role domain.Role, userId string, number int64, purpose domain.AuthPurpose,
	expire time.Time,
) (string, error) {
	secret := s.generateSecret(role, userId, number, purpose)
	tokenToAdd := models.Token{
		UserId:    userId,
		Number:    number,
		Purpose:   int(purpose),
		Secret:    secret,
		ExpiresAt: expire.UnixNano() / 1e+6,
	}
	if _, e := s.repo.AddTokenTX(ctx, tx, role, tokenToAdd); e != nil {
		return "", s.log.Error(e, service_errors.DatabaseErrorRaw)
	}
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"purpose": purpose,
		"secret":  secret,
		"exp":     expire.Unix(),
		"number":  number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, e := token.SignedString([]byte(s.jwtKey))
	if e != nil {
		return "", s.log.Error(e, errors.AuthParseTokenRaw)
	}
	return res, nil
}

func (s *service) generateSecret(role domain.Role, userId string, number int64, purpose domain.AuthPurpose) string {
	toHashElems := []string{
		fmt.Sprintf("%d", role),
		fmt.Sprintf("%d", userId),
		fmt.Sprintf("%d", number),
		fmt.Sprintf("%d", purpose),
		fmt.Sprintf("%d", s.timeAdapter.Now().UnixNano()),
		s.randomAdapter.RandomString(20),
	}
	toHash := strings.Join(toHashElems, "_")
	hash := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hash[:])
}
