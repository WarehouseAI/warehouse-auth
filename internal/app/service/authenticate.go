package service

import (
	m "auth-service/internal/app/model"
	s "auth-service/internal/app/repository/session"
	e "auth-service/internal/pkg/errors/http"
	"context"
)

func Authenticate(sessionId string, sessionRepo s.Repository) (*string, *m.Session, error) {
	userId, newSession, err := sessionRepo.Update(context.Background(), sessionId)

	if err != nil {
		return nil, nil, e.NewHttpErrorByDbStatus(err)
	}

	return userId, newSession, nil
}
