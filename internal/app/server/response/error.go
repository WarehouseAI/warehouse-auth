package response

import (
	e "auth-service/internal/pkg/errors/http"
	"net/http"
)

func WithError(w http.ResponseWriter, httpError error) {
	err := httpError.(e.HttpError)

	JSON(w, err.Status, map[string]interface{}{
		"status":  err.Status,
		"message": err.Err.Error(),
		"trace":   err.Trace,
	})
}
