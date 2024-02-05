package verificationtoken

import m "auth-service/internal/app/model"

type VerificationTokenInterface interface {
	Create(newVerificationToken *m.VerificationToken, email string) error
	Get(condition map[string]interface{}) (*m.VerificationToken, error)
	Delete(condition map[string]interface{}) error
}
