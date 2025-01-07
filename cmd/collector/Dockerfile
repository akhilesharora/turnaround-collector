# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o collector ./cmd/collector

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/collector .

ENTRYPOINT ["./collector"]
