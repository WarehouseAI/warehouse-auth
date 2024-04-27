package reset_token

import (
	"auth-service/internal/repository/models"
	"auth-service/internal/repository/operations/transactions"
	"context"
)

type Repository interface {
	Create(ctx context.Context, tx transactions.Transaction, rt models.ResetToken) (models.ResetToken, error)
	GetById(ctx context.Context, tx transactions.Transaction, id string) (models.ResetToken, error)
	DeleteById(ctx context.Context, tx transactions.Transaction, id string) error
}
