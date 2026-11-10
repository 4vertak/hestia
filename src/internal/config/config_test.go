package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvironmentOverridesDotEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("DB_PASSWORD='file-$#:@ value'\nDB_PORT=5433\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DB_PASSWORD", "environment-value")
	t.Setenv("DB_PORT", "6543")
	cfg, err := load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBPassword != "environment-value" || cfg.DBPort != "6543" {
		t.Fatal("environment must take precedence over .env")
	}
}

func TestDotEnvKeepsQuotedPasswordWithoutChangingEnvironment(t *testing.T) {
	old, exists := os.LookupEnv("DB_PASSWORD")
	if err := os.Unsetenv("DB_PASSWORD"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if exists {
			os.Setenv("DB_PASSWORD", old)
		} else {
			os.Unsetenv("DB_PASSWORD")
		}
	})
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("DB_PASSWORD='file-$#:@ value'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBPassword != "file-$#:@ value" {
		t.Fatal("quoted password was changed")
	}
	if _, exists := os.LookupEnv("DB_PASSWORD"); exists {
		t.Fatal("loading configuration must not mutate the process environment")
	}
}

func TestMissingPasswordFailsWithoutDotEnv(t *testing.T) {
	t.Setenv("DB_PASSWORD", "")
	_, err := load(filepath.Join(t.TempDir(), "missing.env"))
	if err == nil || !strings.Contains(err.Error(), "DB_PASSWORD") {
		t.Fatal("expected a clear error for a missing password")
	}
}

func TestMalformedDotEnvDoesNotExposeItsValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	const sensitive = "sensitive-value-for-parser-test"
	if err := os.WriteFile(path, []byte("DB_PASSWORD='"+sensitive), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := load(path)
	if err == nil || strings.Contains(err.Error(), sensitive) {
		t.Fatal("expected a configuration error without the file contents")
	}
}
