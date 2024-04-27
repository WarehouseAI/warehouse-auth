package rep_converters

import (
	"auth-service/internal/domain"
	wh_converters "auth-service/internal/pkg/utils/converters"
	"auth-service/internal/repository/models"
)

func DomainVerificationToken2ModelVerificationToken(t domain.VerificationTokenInfo) models.VerificationToken {
	return models.VerificationToken{
		UserId:    wh_converters.FastConvertToXid(t.UserId),
		Token:     t.Token,
		SendTo:    t.SendTo,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}
}

func ModelVerificationToken2DomainVerificationToken(t models.VerificationToken) domain.VerificationTokenInfo {
	return domain.VerificationTokenInfo{
		ID:        t.ID.String(),
		UserId:    t.ID.String(),
		Token:     t.Token,
		SendTo:    t.SendTo,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}
}

func DomainResetToken2ModelResetToken(t domain.ResetTokenInfo) models.ResetToken {
	return models.ResetToken{
		UserId:    wh_converters.FastConvertToXid(t.UserId),
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}
}

func ModelResetToken2DomainResetToken(t models.ResetToken) domain.ResetTokenInfo {
	return domain.ResetTokenInfo{
		ID:        t.ID.String(),
		UserId:    t.UserId.String(),
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}
}
