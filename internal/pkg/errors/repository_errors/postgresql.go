package repository_errors

import "errors"

var (
	PostgresqlScanRaw            = errors.New("postgresql scan error")
	PostgresqlQueryRaw           = errors.New("postgresql query error")
	PostgresqlGetRaw             = errors.New("postgresql get error")
	PostgresqlQueryRowRaw        = errors.New("postgresql query row error")
	PostgresqlExecRaw            = errors.New("postgresql exec error")
	PostgresqlRowsAffectedRaw    = errors.New("postgresql rows affected error")
	PostgresqlLastInsertIdRaw    = errors.New("postgresql last insert id error")
	PostgresqlRebindRaw          = errors.New("postgresql rebind error")
	PostgresqlRowsCloseRaw       = errors.New("postgresql rows close error")
	PostgresqlNoRowsWereAffected = errors.New("no rows were affected")

	PgTx = errors.New("postgresql transaction error") // PostgresqlTransaction
)
