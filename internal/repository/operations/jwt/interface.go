package jwt

import (
	"context"

	"github.com/warehouse/auth-service/internal/domain"
	"github.com/warehouse/auth-service/internal/repository/models"
	"github.com/warehouse/auth-service/internal/repository/operations/transactions"
)

type Repository interface {
	DropAllTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) error
	DropTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string, number int64) error
	FindNumberTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) (int64, error)
	AddTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token models.Token) (models.Token, error)
	CheckTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token models.Token) (models.Token, error)
	DropOldTokens(ctx context.Context, tx transactions.Transaction, timestamp int64) error

	GetTokenMap() map[domain.Role]string
}
