package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Grpc struct {
		UserHost string
		AuthHost string
		Post     int
		UserAddr string
		AuthAddr string
	}

	Auth struct {
		Key                 string
		AccessTokenTimeout  time.Duration
		RefreshTokenTimeout time.Duration
		AuthTimeout         time.Duration
	}

	Timeouts struct {
		AuthTimeout    time.Duration
		RequestTimeout time.Duration
		AccCookie      time.Duration
	}

	Postgres struct {
		DSN     string
		CertLoc string
	}

	Rabbit struct {
		URL       string
		MailQueue string
		UserQueue string
	}

	Server struct {
		Mode           string
		Port           int
		AllowedOrigins []string
	}

	Time struct {
		Locale int64
	}

	Config struct {
		Server   Server
		Rabbit   Rabbit
		Auth     Auth
		Timeouts Timeouts
		Postgres Postgres
		Grpc     Grpc
		Time     Time
	}
)

func loadJwtKey(v *viper.Viper, isProd bool) (string, error) {
	if isProd {
		path := v.GetString("apis.jwt")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		var jwtCreds struct {
			Key string `json:"key"`
		}
		if err := json.Unmarshal(data, &jwtCreds); err != nil {
			return "", err
		}

		return jwtCreds.Key, nil
	} else {
		value := v.GetString("apis.jwt")
		return value, nil
	}
}

func loadPgSource(v *viper.Viper, isProd bool) (string, error) {
	if isProd {
		path := v.GetString("pgSource")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		var pgSource struct {
			Source string `json:"source"`
		}
		if err := json.Unmarshal(data, &pgSource); err != nil {
			return "", err
		}

		return pgSource.Source, nil
	} else {
		pgSource := v.GetString("pgSource")
		return pgSource, nil
	}
}

func generateRabbitUrl(v *viper.Viper) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		v.GetString("rabbitmq.user"),
		v.GetString("rabbitmq.password"),
		v.GetString("rabbitmq.host"),
		v.GetInt("rabbitmq.port"),
	)
}

func Timeout(v *viper.Viper, kind string) time.Duration {
	return time.Second * time.Duration(v.GetInt(fmt.Sprintf("timeouts.%s", kind)))
}

func NewConfig(cfgPath string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(cfgPath)
	v.SetConfigName("config")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	mode := v.GetString("mode")

	jwtKey, err := loadJwtKey(v, mode == "prod")
	if err != nil {
		return nil, err
	}

	pgSource, err := loadPgSource(v, mode == "prod")
	if err != nil {
		return nil, err
	}

	return &Config{
		Auth: Auth{
			Key:                 jwtKey,
			AccessTokenTimeout:  Timeout(v, "access_token"),  // таймаут цифрами для ttl токена
			RefreshTokenTimeout: Timeout(v, "refresh_token"), // таймаут цифрами для ttl токена
			AuthTimeout:         Timeout(v, "request"),
		},
		Timeouts: Timeouts{
			RequestTimeout: v.GetDuration("request_timeout.request"), // общие таймауты (можно переносить между сервисами)
			AuthTimeout:    v.GetDuration("request_timeout.auth"),
			AccCookie:      v.GetDuration("acc_cookie"),
		},
		Grpc: Grpc{
			UserHost: v.GetString("grpc.user_host"),
			AuthHost: v.GetString("grpc.auth_host"),
			Post:     v.GetInt("grpc.port"),
			AuthAddr: v.GetString("grpc.auth_addr"),
			UserAddr: v.GetString("grpc.user_addr"),
		},
		Postgres: Postgres{
			DSN:     pgSource,
			CertLoc: v.GetString("pgCertLoc"),
		},
		Server: Server{
			Mode:           mode,
			Port:           v.GetInt("server.port"),
			AllowedOrigins: v.GetStringSlice("server.allowed_origins"),
		},

		Rabbit: Rabbit{
			URL:       generateRabbitUrl(v),
			MailQueue: v.GetString("rabbitmq.queues.mail"),
			UserQueue: v.GetString("rabbitmq.queues.user"),
		},

		Time: Time{
			Locale: v.GetInt64("locale"),
		},
	}, nil

}
