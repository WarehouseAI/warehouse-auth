package verification_token

import (
	"context"

	"github.com/warehouse/auth-service/internal/repository/models"
	"github.com/warehouse/auth-service/internal/repository/operations/transactions"
)

type Repository interface {
	Create(ctx context.Context, tx transactions.Transaction, vt models.VerificationToken) (models.VerificationToken, error)
	GetById(ctx context.Context, tx transactions.Transaction, id string) (models.VerificationToken, error)
	DeleteById(ctx context.Context, tx transactions.Transaction, id string) error
}
