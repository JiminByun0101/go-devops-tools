# Dockerfile
# Stage 1: Build the Go application

FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /k8s-watchdog ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /k8s-watchdog .

CMD ["/app/k8s-watchdog"]