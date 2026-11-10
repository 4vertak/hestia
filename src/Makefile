.DEFAULT_GOAL := build
.PHONY: run run-local run-docker docker-down build swagger swagger-check migrate-up migrate-down test vet clean

GO ?= go
SWAG_VERSION := v1.16.6
SWAG = $(GO) run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

# Generated docs are tracked. Regenerate explicitly after changing annotations.
docs/docs.go:
	$(MAKE) swagger

swagger:
	$(SWAG) init -g cmd/server/main.go --output docs --packageName docs

swagger-check:
	@set -eu; generated=$$(mktemp -d); trap 'rm -rf "$$generated"' EXIT; \
	$(SWAG) init -g cmd/server/main.go --output "$$generated" --packageName docs; \
	diff -ru docs "$$generated"

build: docs/docs.go
	$(GO) build -o bin/server ./cmd/server

test: docs/docs.go
	$(GO) test -count=1 -timeout 30s ./...

vet: docs/docs.go
	$(GO) vet ./...

run: run-local

run-local: build
	./bin/server

run-docker:
	docker compose up --build

docker-down:
	docker compose down

# Share the application's .env loading and embedded migrations; no migrate CLI.
migrate-up:
	$(GO) run ./cmd/migrate up

migrate-down:
	$(GO) run ./cmd/migrate down

clean:
	rm -rf bin
