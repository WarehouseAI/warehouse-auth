package dataservice

import (
	"auth-service/configs"
	d "auth-service/internal/app/dataservice/operations"
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

func NewResetTokenDatabase(config configs.Config) *d.TokenDatabase[m.ResetToken] {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Database,
		config.Postgres.Port,
	)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	return &d.TokenDatabase[m.ResetToken]{DB: db}
}

func NewVerificationTokenDatabase(config configs.Config) *d.TokenDatabase[m.VerificationToken] {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Database,
		config.Postgres.Port,
	)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	return &d.TokenDatabase[m.VerificationToken]{DB: db}
}
