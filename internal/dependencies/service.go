package dependencies

import (
	"auth-service/internal/service/auth"
	"auth-service/internal/service/jwt"
)

func (d *dependencies) AuthService() auth.Service {
	if d.authService == nil {
		d.authService = auth.NewService(
			*d.cfg,
			d.PgxTransactionRepo(),
			d.JwtRepo(),
			d.TimeAdapter(),
			d.JwtService(),
			d.log,
			d.UserAdapter(),
			d.MailAdapter(),
			d.VerificationTokenRepo(),
			d.ResetTokenRepo(),
		)
	}

	return d.authService
}

func (d *dependencies) JwtService() jwt.Service {
	if d.jwtService == nil {
		d.jwtService = jwt.NewService(
			d.log,
			d.PgxTransactionRepo(),
			d.JwtRepo(),
			d.cfg.Auth,
			d.TimeAdapter(),
			d.RandomAdapter(),
		)
	}

	return d.jwtService
}
