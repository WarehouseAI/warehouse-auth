package reset_token

import (
	"auth-service/internal/pkg/errors/repository_errors"
	"auth-service/internal/repository/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (r *repositoryPG) getResetTokenByCondition(
	ctx context.Context,
	executor sqlx.ExtContext,
	condition string,
	params ...interface{},
) ([]models.ResetToken, error) {
	query := `
    SELECT rt.id, rt.user_id, rt.token, rt.created_at, rt.expires_at
    FROM reset_tokens as rt
  `
	query = fmt.Sprintf("%s %s", query, condition)

	var list []models.ResetToken
	err := sqlx.SelectContext(ctx, executor, &list, query, params...)
	if err != nil {
		return []models.ResetToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlGetRaw, query)
	}

	return list, nil
}
