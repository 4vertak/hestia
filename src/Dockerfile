FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /server .

EXPOSE 8080
CMD ["./server"]
