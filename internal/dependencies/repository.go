package dependencies

import (
	"auth-service/internal/repository/operations/jwt"
	"auth-service/internal/repository/operations/reset_token"
	"auth-service/internal/repository/operations/transactions"
	"auth-service/internal/repository/operations/verification_token"
)

func (d *dependencies) PgxTransactionRepo() transactions.Repository {
	if d.pgxTransactionRepo == nil {
		d.pgxTransactionRepo = transactions.NewPgxRepository(d.PostgresClient())
	}
	return d.pgxTransactionRepo
}

func (d *dependencies) JwtRepo() jwt.Repository {
	if d.jwtRepo == nil {
		d.jwtRepo = jwt.NewPGRepository(d.log, d.PostgresClient())
	}
	return d.jwtRepo
}

func (d *dependencies) VerificationTokenRepo() verification_token.Repository {
	if d.verificationTokenRepo == nil {
		d.verificationTokenRepo = verification_token.NewPGRepository(d.log, d.PostgresClient())
	}

	return d.verificationTokenRepo
}

func (d *dependencies) ResetTokenRepo() reset_token.Repository {
	if d.resetTokenRepo == nil {
		d.resetTokenRepo = reset_token.NewPGRepository(d.log, d.PostgresClient())
	}

	return d.resetTokenRepo
}
