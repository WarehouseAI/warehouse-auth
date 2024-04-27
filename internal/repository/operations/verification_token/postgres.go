package verification_token

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
		log: log.Named("pg_verification_tokens"),
	}
}

func (r *repositoryPG) Create(
	ctx context.Context,
	tx transactions.Transaction,
	vt models.VerificationToken,
) (models.VerificationToken, error) {
	query := `
    INSERT INTO verification_tokens (user_id, token, send_to)
    VALUES(:user_id, :token, :send_to)
  `

	res, err := tx.Txm().NamedExecContext(ctx, query, vt)
	if err != nil {
		return models.VerificationToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.VerificationToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	if rowsAffected != 1 {
		return models.VerificationToken{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	return vt, nil
}

func (r *repositoryPG) GetById(
	ctx context.Context,
	tx transactions.Transaction,
	id string,
) (models.VerificationToken, error) {
	cond := `WHERE vt.id = $1`
	list, err := r.getVerificationTokenByCondition(ctx, tx.Txm(), cond, id)
	if err != nil {
		return models.VerificationToken{}, err
	}

	return list[0], nil
}

func (r *repositoryPG) DeleteById(
	ctx context.Context,
	tx transactions.Transaction,
	id string,
) error {
	query := `DELETE FROM verification_tokens WHERE id=$1`
	_, err := tx.Txm().ExecContext(ctx, query, id)
	if err != nil {
		return r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}
	return nil
}
