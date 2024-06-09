package reset_token

import (
	"context"

	"github.com/warehouse/auth-service/internal/repository/models"
	"github.com/warehouse/auth-service/internal/repository/operations/transactions"
)

type Repository interface {
	Create(ctx context.Context, tx transactions.Transaction, rt models.ResetToken) (models.ResetToken, error)
	GetById(ctx context.Context, tx transactions.Transaction, id string) (models.ResetToken, error)
	DeleteById(ctx context.Context, tx transactions.Transaction, id string) error
}
