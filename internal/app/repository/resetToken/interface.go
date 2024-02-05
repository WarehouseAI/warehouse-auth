package resettoken

import m "auth-service/internal/app/model"

type Repository interface {
	Create(newResetToken *m.ResetToken) error
	Get(condition map[string]interface{}) (*m.ResetToken, error)
	Delete(condition map[string]interface{}) error
}
