package dataservice

import (
	m "auth-service/internal/app/model"
	"context"
	"mime/multipart"
)

type ResetTokenInterface interface {
	Create(newResetToken *m.ResetToken) error
	Get(condition map[string]interface{}) (*m.ResetToken, error)
	Delete(condition map[string]interface{}) error
}

type VerificationTokenInterface interface {
	Create(newVerificationToken *m.VerificationToken) error
	Get(condition map[string]interface{}) (*m.VerificationToken, error)
	Delete(condition map[string]interface{}) error
}

type SessionInterface interface {
	Create(ctx context.Context, userId string) (*m.Session, error)
	Get(ctx context.Context, sessionId string) (*m.Session, error)
	Delete(ctx context.Context, sessionId string) error
	Update(ctx context.Context, sessionId string) (*string, *m.Session, error)
}

type PictureInterface interface {
	UploadFile(file multipart.File, fileName string) (string, error)
	DeleteImage(fileName string) error
}
