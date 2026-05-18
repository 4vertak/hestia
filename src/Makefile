.PHONY: run build swagger migrate-up migrate-down test

include .env
export

run:
	docker-compose up --build

run-local:
	export $$(cat .env | xargs) && go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

swagger:
	swag init -g cmd/server/main.go --output docs

migrate-up:
	migrate -path internal/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path internal/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

test:
	export PATH=$$PATH:$$(go env GOPATH)/bin && \
	swag init -g cmd/server/main.go --output docs && \
	go test -count=1 ./... -timeout 30s -v 2>&1 | sed '/no test files/d'

clean:
	rm -rf ./docs
	rm -rf ./bin