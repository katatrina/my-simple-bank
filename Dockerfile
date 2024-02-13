# Build stage
FROM golang:1.22.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.12.0
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY db/migration ./migration
COPY start-migration.sh .
COPY wait-for .

EXPOSE 8080

ENTRYPOINT ["/app/start-migration.sh"]
CMD ["/app/main"]