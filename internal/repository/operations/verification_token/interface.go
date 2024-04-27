package verification_token

import (
	"auth-service/internal/repository/models"
	"auth-service/internal/repository/operations/transactions"
	"context"
)

type Repository interface {
	Create(ctx context.Context, tx transactions.Transaction, vt models.VerificationToken) (models.VerificationToken, error)
	GetById(ctx context.Context, tx transactions.Transaction, id string) (models.VerificationToken, error)
	DeleteById(ctx context.Context, tx transactions.Transaction, id string) error
}
