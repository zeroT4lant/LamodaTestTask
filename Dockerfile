FROM golang:alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main app/cmd/main.go

FROM alpine as runner

WORKDIR /app

COPY .env /app/.env

COPY --from=builder /src/main /app/


ENTRYPOINT ./main

CMD ["/app"]