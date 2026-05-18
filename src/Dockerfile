FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN swag init -g cmd/server/main.go --output docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /server .

COPY --from=builder /app/internal/database/migrations/ /root/migrations
COPY --from=builder /app/docs/ /root/docs/

EXPOSE 8080
CMD ["./server"]