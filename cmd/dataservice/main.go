package dataservice

import (
	"auth-service/configs"
	"auth-service/internal/app/model"
	rt "auth-service/internal/app/repository/resetToken"
	s "auth-service/internal/app/repository/session"
	vt "auth-service/internal/app/repository/verificationToken"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSessionDatabase(config configs.Config) *s.Database {
	rClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})

	return &s.Database{
		DB: rClient,
	}
}

func NewSqlDatabase(config configs.Config) (*rt.Database, *vt.Database) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Database,
	)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.Raw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&model.VerificationToken{})
	db.AutoMigrate(&model.ResetToken{})

	return &rt.Database{DB: db}, &vt.Database{DB: db}
}
