package database

import (
	"net/url"
	"testing"

	"github.com/4vertak/hestia/internal/config"
)

func TestPostgresDSNPreservesCredentialsAndIPv6(t *testing.T) {
	cfg := &config.Config{
		DBUser:     "user@example",
		DBPassword: "synthetic-$#:@ /? password",
		DBHost:     "::1",
		DBPort:     "5433",
		DBName:     "store",
		DBSSLMode:  "disable",
	}
	dsn, err := url.Parse(postgresDSN(cfg))
	if err != nil {
		t.Fatal("invalid connection URL")
	}
	password, _ := dsn.User.Password()
	if dsn.User.Username() != cfg.DBUser || password != cfg.DBPassword {
		t.Fatal("connection URL changed the credentials")
	}
	if dsn.Hostname() != cfg.DBHost || dsn.Port() != cfg.DBPort || dsn.Query().Get("sslmode") != cfg.DBSSLMode {
		t.Fatal("connection URL changed the host, port or SSL mode")
	}
}
