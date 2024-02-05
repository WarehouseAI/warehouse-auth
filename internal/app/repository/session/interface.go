package session

import (
	m "auth-service/internal/app/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, userId string) (*m.Session, error)
	Get(ctx context.Context, sessionId string) (*m.Session, error)
	Delete(ctx context.Context, sessionId string) error
	Update(ctx context.Context, sessionId string) (*string, *m.Session, error)
}
