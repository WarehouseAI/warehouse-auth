package dependencies

import "github.com/warehouse/auth-service/internal/handler/middlewares"

func (d *dependencies) HandlerMiddleware() middlewares.Middleware {
	if d.handlerMiddleware == nil {
		d.handlerMiddleware = middlewares.NewMiddleware(
			d.log,
			d.cfg.Timeouts,
			d.JwtService(),
		)
	}

	return d.handlerMiddleware
}
