package dependencies

import "auth-service/internal/handler/http"

func (d *dependencies) AuthHandler() http.Handler {
	if d.authHandler == nil {
		d.authHandler = http.NewAuthHandler(
			d.cfg.Server,
			d.cfg.Timeouts,
			d.JwtService(),
			d.AuthService(),
			d.TimeAdapter(),
			d.UserAdapter(),
			d.WarehouseJsonRequestHandler(),
			d.HandlerMiddleware(),
		)
	}

	return d.authHandler
}
