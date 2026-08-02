FROM golang:1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ordersystem ./cmd/ordersystem

FROM alpine:latest
WORKDIR /app
COPY --from=builder /ordersystem .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000 50051 8080

ENTRYPOINT ["/app/ordersystem"]
