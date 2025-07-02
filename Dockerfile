# Builder layer
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o packaging-service ./cmd/server/main.go

# Execution layer
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/packaging-service .

EXPOSE 50051

CMD ["./packaging-service"]
