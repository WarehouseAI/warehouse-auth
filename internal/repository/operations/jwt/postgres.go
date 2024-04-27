package jwt

import (
	"auth-service/internal/db"
	"auth-service/internal/domain"
	"auth-service/internal/pkg/errors"
	"auth-service/internal/pkg/errors/repository_errors"
	"auth-service/internal/pkg/logger"
	"auth-service/internal/repository/models"
	"auth-service/internal/repository/operations/transactions"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repositoryPG struct {
	log      logger.Logger
	pg       *db.PostgresClient
	tokenMap map[domain.Role]string
}

func NewPGRepository(
	log logger.Logger,
	client *db.PostgresClient,
) Repository {
	return &repositoryPG{
		log:      log.Named("pg_jwt_repo"),
		pg:       client,
		tokenMap: UserTokenMap,
	}
}

func (r *repositoryPG) FindNumberTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) (int64, error) {
	var numbers []int64
	queryString := `
		SELECT  number 
		FROM %s 
		WHERE user_id=$1 AND purpose=0 
		ORDER BY number
	`
	query := fmt.Sprintf(queryString, r.tokenMap[role])
	err := tx.Txm().SelectContext(ctx, &numbers, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			numbers = []int64{}
		}
		return 0, r.log.ErrorRepo(err, repository_errors.PostgresqlQueryRowRaw, query)
	}

	return r.findNumbers(numbers)
}

func (r *repositoryPG) AddTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token models.Token) (models.Token, error) {
	queryString := `
		INSERT INTO %s (user_id, number, purpose, secret, expires_at)
		VALUES(:user_id, :number, :purpose, :secret, :expires_at)
	`
	query := fmt.Sprintf(queryString, r.tokenMap[role])

	res, err := tx.Txm().NamedExecContext(ctx, query, token)
	if err != nil {
		return models.Token{}, r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Token{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	if rowsAffected != 1 {
		return models.Token{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	return token, err
}

func (r *repositoryPG) DropTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string, number int64) error {
	queryString := `DELETE FROM %s WHERE user_id=$1 AND number=$2`
	query := fmt.Sprintf(queryString, r.tokenMap[role])
	_, err := tx.Txm().ExecContext(ctx, query, userId, number)
	if err != nil {
		return r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}

	return nil
}

func (r *repositoryPG) DropAllTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, userId string) error {
	queryString := `DELETE FROM %s WHERE user_id=$1`
	query := fmt.Sprintf(queryString, r.tokenMap[role])
	_, err := tx.Txm().ExecContext(ctx, query, userId)
	if err != nil {
		return r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}
	return nil
}

func (r *repositoryPG) CheckTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token models.Token) (models.Token, error) {
	queryString := `
		SELECT id, number, purpose, secret, expires_at
		FROM %s
		WHERE user_id=:user_id
			AND number=:number
			AND purpose=:purpose
			AND secret=:secret
	`
	query := fmt.Sprintf(queryString, r.tokenMap[role])

	rows, err := sqlx.NamedQueryContext(ctx, tx.Txm(), query, token)
	if err != nil {
		return models.Token{}, r.log.ErrorRepo(err, repository_errors.PostgresqlGetRaw, query)
	}

	if !rows.Next() {
		return models.Token{}, errors.TokenDoesNotExist
	}

	err = rows.StructScan(&token)
	if err != nil {
		return models.Token{}, r.log.ErrorRepo(err, repository_errors.PostgresqlScanRaw, query)
	}

	return token, nil
}

func (r *repositoryPG) DropOldTokens(ctx context.Context, tx transactions.Transaction, timestamp int64) error {
	for _, t := range r.tokenMap {
		query := `
			DELETE FROM %s WHERE number IN (SELECT number FROM %s WHERE purpose=$1 AND expires_at<=$2)
		`
		_, err := tx.Txm().ExecContext(ctx, fmt.Sprintf(query, t, t), 1, timestamp)
		if err != nil {
			return r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
		}
	}
	return nil
}

func (r *repositoryPG) GetTokenMap() map[domain.Role]string {
	return r.tokenMap
}

func (r *repositoryPG) findNumbers(numbers []int64) (int64, error) {

	if len(numbers) == 0 {
		return 0, nil
	}
	if numbers[len(numbers)-1] == int64(len(numbers)-1) {
		return int64(len(numbers)), nil
	}
	for i, n := range numbers {
		if n != int64(i) {
			return int64(i), nil
		}
	}

	return 0, errors.AuthNumberAssignmentFailedRaw
}
