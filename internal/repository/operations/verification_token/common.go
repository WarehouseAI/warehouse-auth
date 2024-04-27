package verification_token

import (
	"auth-service/internal/pkg/errors/repository_errors"
	"auth-service/internal/repository/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (r *repositoryPG) getVerificationTokenByCondition(
	ctx context.Context,
	executor sqlx.ExtContext,
	condition string,
	params ...interface{},
) ([]models.VerificationToken, error) {
	query := `
    SELECT vt.id, vt.user_id, vt.token, vt.send_to, vt.created_at, vt.expires_at
    FROM verification_tokens as vt
  `
	query = fmt.Sprintf("%s %s", query, condition)

	var list []models.VerificationToken
	err := sqlx.SelectContext(ctx, executor, &list, query, params...)
	if err != nil {
		return []models.VerificationToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlGetRaw, query)
	}

	return list, nil
}
