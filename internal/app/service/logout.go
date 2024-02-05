package service

import (
	"auth-service/internal/app/repository/session"
	e "auth-service/internal/pkg/errors/http"
	"context"
)

func Logout(sessionId string, session session.Repository) error {
	if err := session.Delete(context.Background(), sessionId); err != nil {
		return e.NewHttpErrorByDbStatus(err)
	}

	return nil
}
