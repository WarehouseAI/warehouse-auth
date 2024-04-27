package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresClient struct {
	DB *sqlx.DB
}

func NewPostgresClient(DSN string, pgCertLoc string) (*PostgresClient, error) {
	var client *sqlx.DB
	if pgCertLoc != "" && !strings.Contains(DSN, "disable") {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(pgCertLoc)
		if err != nil {
			return nil, fmt.Errorf("error reading cert: %w", err)
		}
		rootCertPool.AppendCertsFromPEM(pem)
		connCfg, err := pgx.ParseConfig(DSN)
		if err != nil {
			return nil, fmt.Errorf("error parsing postgres cdn: %w", err)
		}
		connCfg.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
		db := stdlib.OpenDB(*connCfg)
		client = sqlx.NewDb(db, "pgx")
	} else {
		var err error
		client, err = sqlx.Connect("pgx", DSN)
		if err != nil {
			return nil, fmt.Errorf("error while connecting to postgres %w", err)
		}
	}
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("error while ping to postgres %w", err)
	}

	client.SetMaxOpenConns(100)
	client.SetMaxIdleConns(100)
	client.SetConnMaxLifetime(10 * time.Second)

	return &PostgresClient{
		DB: client,
	}, nil
}

func (client *PostgresClient) Close() error {
	return client.DB.Close()
}
