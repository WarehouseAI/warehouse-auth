package dependencies

import (
	"github.com/warehouse/auth-service/internal/handler/grpc"
	"github.com/warehouse/auth-service/internal/handler/http"
)

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

func (d *dependencies) AuthGrpcHandler() *grpc.AuthHandler {
	if d.authGrpcHandler == nil {
		d.authGrpcHandler = grpc.NewAuthHandler(
			d.cfg.Timeouts,
			d.log,
			d.JwtService(),
		)
	}

	return d.authGrpcHandler
}
