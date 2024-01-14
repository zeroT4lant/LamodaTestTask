# Build Stage
#FROM golang:1.18.6-alpine3.15 AS builder
FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go build -o main ./cmd
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY configs/app.env .
COPY start.sh .
COPY wait-for.sh .
COPY migrations ./migrations

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]