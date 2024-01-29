package sessiondata

import (
	m "auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	DB *redis.Client
}

// TODO: Подумать над errorHandler как в sql функциях

func (d *Database) Create(ctx context.Context, userId string) (*m.Session, error) {
	TTL := 24 * time.Hour
	sessionId := uuid.Must(uuid.NewV4()).String()

	sessionPayload := m.SessionPayload{
		UserId:    userId,
		CreatedAt: time.Now(),
	}

	marshaledPayload, err := json.Marshal(sessionPayload)

	if err != nil {
		return nil, e.NewDBError(e.DbSystem, err.Error(), fmt.Errorf("Something went wrong."))
	}

	if err := d.DB.Set(ctx, sessionId, marshaledPayload, TTL).Err(); err != nil {
		return nil, e.NewDBError(e.DbSystem, err.Error(), fmt.Errorf("Something went wrong."))
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTL}, nil
}

func (d *Database) Get(ctx context.Context, sessionId string) (*m.Session, error) {
	var sessionPayload m.SessionPayload

	record := d.DB.Get(ctx, sessionId)
	recordTTL := d.DB.TTL(ctx, sessionId)

	if record.Err() != nil {
		return nil, e.NewDBError(e.DbNotFound, record.Err().Error(), fmt.Errorf("Session not found"))
	}

	recordInfo, _ := record.Result()
	TTLInfo, _ := recordTTL.Result()

	if err := json.Unmarshal([]byte(recordInfo), &sessionPayload); err != nil {
		return nil, e.NewDBError(e.DbSystem, err.Error(), fmt.Errorf("Something went wrong."))
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTLInfo}, nil
}

func (d *Database) Delete(ctx context.Context, sessionId string) error {
	if err := d.DB.Del(ctx, sessionId).Err(); err != nil {
		return e.NewDBError(e.DbSystem, err.Error(), fmt.Errorf("Something went wrong."))
	}

	return nil
}

func (d *Database) Update(ctx context.Context, sessionId string) (*string, *m.Session, error) {
	session, err := d.Get(ctx, sessionId)

	if err != nil {
		return nil, nil, err
	}

	if err := d.Delete(ctx, sessionId); err != nil {
		return nil, nil, err
	}

	newSession, err := d.Create(ctx, session.Payload.UserId)

	return &session.Payload.UserId, newSession, nil
}
