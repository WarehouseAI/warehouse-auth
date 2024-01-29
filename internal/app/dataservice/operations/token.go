package operations

import (
	m "auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/db"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type TokenDatabase[T m.Tokens] struct {
	DB *gorm.DB
}

func (d *TokenDatabase[T]) errorHandle(err error) error {
	if err == nil {
		return nil
	}

	// Добавлять новые ошибки в этот свитч и использовать потом внутри if с ошибкой
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		switch pgErr.Code {
		case "23505":
			return e.NewDBError(e.DbExist, err.Error(), fmt.Errorf("Token for this user is already exists."))

		case "20000":
			return e.NewDBError(e.DbNotFound, err.Error(), fmt.Errorf("Token not found."))

		case "25503":
			return e.NewDBError(e.DbForeignKey, err.Error(), fmt.Errorf("Invalid foreign key in token declaration."))
		}
	}

	return e.NewDBError(e.DbSystem, err.Error(), fmt.Errorf("Something went wrong."))
}

func (d *TokenDatabase[T]) Create(token *T) error {
	if err := d.DB.Create(token).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}

func (d *TokenDatabase[T]) Get(conditions map[string]interface{}) (*T, error) {
	var token T

	if err := d.DB.Where(conditions).First(&token).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &token, nil
}

func (d *TokenDatabase[T]) Delete(condition map[string]interface{}) error {
	var token T

	if err := d.DB.Where(condition).Delete(&token).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}
