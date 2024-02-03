package dataservice

import (
	"auth-service/configs"
	d "auth-service/internal/app/dataservice/operations"
	"auth-service/internal/app/model"
	m "auth-service/internal/app/model"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSessionDatabase(config configs.Config) *d.SessionDatabase {
	rClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})

	return &d.SessionDatabase{
		DB: rClient,
	}
}

func NewSqlDatabase(config configs.Config) (*d.TokenDatabase[m.ResetToken], *d.TokenDatabase[m.VerificationToken]) {
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

	return &d.TokenDatabase[m.ResetToken]{DB: db}, &d.TokenDatabase[m.VerificationToken]{DB: db}
}
