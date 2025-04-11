FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY  . .
RUN go build -o /app/bin/hermes main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/hermes /app/hermes

CMD ["/app/hermes"]
