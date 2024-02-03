package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Env            string   `envconfig:"SERVER_ENV"`
		Port           string   `envconfig:"SERVER_PORT"`
		GrpcPort       string   `envconfig:"SERVER_GRPC_PORT"`
		GrpcHost       string   `envconfig:"SERVER_GRPC_HOST"`
		UserAddr       string   `envconfig:"SERVER_USER_ADDR"`
		AllowedOrigins []string `endconfig:"SERVER_ALLOWED_ORIGINS"`
	}

	Postgres struct {
		Host     string `envconfig:"PSQL_DB_HOST"`
		Port     int    `envconfig:"PSQL_DB_PORT"`
		Database string `envconfig:"PSQL_DB_NAME"`
		Username string `envconfig:"PSQL_DB_USER"`
		Password string `envconfig:"PSQL_DB_PASS"`
	}

	Redis struct {
		Host     string `envconfig:"REDIS_DB_HOST"`
		Port     int    `envconfig:"REDIS_DB_PORT"`
		Password string `envconfig:"REDIS_DB_PASSWORD"`
	}

	Rabbit struct {
		Host     string   `envconfig:"RMQ_HOST"`
		Port     int      `envconfig:"RMQ_PORT"`
		Username string   `envconfig:"RMQ_USER"`
		Password string   `envconfig:"RMQ_PASS"`
		Queues   []string `envconfig:"RMQ_QUEUES"`
	}
}

func ReadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, err
}
