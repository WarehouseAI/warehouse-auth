package middlewares

import (
	"auth-service/internal/handler/converters"
	"auth-service/internal/handler/writers"
	"auth-service/internal/pkg/errors"
	"context"
	"net/http"
	"time"
)

func (m *middleware) acquireWorker(ctx context.Context) *errors.Error {
	select {
	case <-ctx.Done():
		return &errors.Error{
			Code:   429,
			Reason: "too many requests",
		}
	case m.queue <- struct{}{}:
		return nil
	}
}

func (m *middleware) releaseWorker() {
	<-m.queue
}

func (m *middleware) QueueMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		e := m.acquireWorker(ctx)
		if e != nil {
			writers.SendJSON(w, int(e.Code), converters.MakeJsonErrorResponseWithErrorsError(e))
			return
		}
		defer m.releaseWorker()
		h.ServeHTTP(w, r)
	})
}