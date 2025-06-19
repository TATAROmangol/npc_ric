FROM golang:1.24-alpine3.21 AS builder 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./forms ./cmd/main.go

FROM alpine:3.21
WORKDIR /app 
COPY --from=builder /app/generator ./bin/
COPY --from=builder /app/internal/storage/migrations ./internal/storage/migrations
CMD ["./bin/generator"]