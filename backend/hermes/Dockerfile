FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.work ./
COPY libs/go.mod ./libs/
COPY hermes/go.mod ./hermes/

RUN go work sync && cd hermes && go mod download

COPY  ./libs ./libs
COPY  ./hermes ./hermes
RUN go build -o /app/bin/hermes ./hermes/cmd/hermes/main.go

FROM alpine:latest
RUN apk add tzdata
COPY --from=builder /app/bin/hermes /app/hermes
COPY hermes/migrations /app/migrations
COPY hermes/scripts/entrypoint.sh /app/entrypoint.sh
COPY hermes/.env /app/.env
RUN chmod +x /app/entrypoint.sh
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

CMD ["/app/entrypoint.sh"]
