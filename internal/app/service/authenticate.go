package service

import (
	"auth-service/internal/app/dataservice"
	m "auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/http"
	"context"
)

func Authenticate(sessionId string, sessionRepo dataservice.SessionInterface) (*string, *m.Session, error) {
	userId, newSession, err := sessionRepo.Update(context.Background(), sessionId)

	if err != nil {
		return nil, nil, e.NewHttpErrorByDbStatus(err)
	}

	return userId, newSession, nil
}
