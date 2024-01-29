package service

import (
	"auth-service/internal/app/dataservice"
	e "auth-service/internal/pkg/errors/http"
	"context"
)

func Logout(sessionId string, session dataservice.SessionInterface) error {
	if err := session.Delete(context.Background(), sessionId); err != nil {
		return e.NewHttpErrorByDbStatus(err)
	}

	return nil
}
