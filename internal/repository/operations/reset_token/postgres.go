package reset_token

import (
	"auth-service/internal/db"
	"auth-service/internal/pkg/errors/repository_errors"
	"auth-service/internal/pkg/logger"
	"auth-service/internal/repository/models"
	"auth-service/internal/repository/operations/transactions"
	"context"
)

type repositoryPG struct {
	log logger.Logger
	pg  *db.PostgresClient
}

func NewPGRepository(log logger.Logger, client *db.PostgresClient) Repository {
	return &repositoryPG{
		pg:  client,
		log: log.Named("pg_reset_tokens"),
	}
}

func (r *repositoryPG) Create(ctx context.Context, tx transactions.Transaction, rt models.ResetToken) (models.ResetToken, error) {
	query := `
    INSERT INTO reset_tokens (user_id, token)
    VALUES(:user_id, :token)
  `

	res, err := tx.Txm().NamedExecContext(ctx, query, rt)
	if err != nil {
		return models.ResetToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.ResetToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	if rowsAffected != 1 {
		return models.ResetToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	return rt, nil
}

func (r *repositoryPG) GetById(ctx context.Context, tx transactions.Transaction, id string) (models.ResetToken, error) {
	cond := `WHERE rt.id = $1`
	list, err := r.getResetTokenByCondition(ctx, tx.Txm(), cond, id)
	if err != nil {
		return models.ResetToken{}, err
	}

	return list[0], nil
}

func (r *repositoryPG) DeleteById(
	ctx context.Context,
	tx transactions.Transaction,
	id string,
) error {
	query := `DELETE FROM reset_tokens WHERE id=$1`
	_, err := tx.Txm().ExecContext(ctx, query, id)
	if err != nil {
		return r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}
	return nil
}
