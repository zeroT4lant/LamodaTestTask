DB_URL=postgresql://postgres:secret@localhost:5431/lamoda_db?sslmode=disable

up:
	docker-compose up --build

test:
	cd app && GO111MODULE=on go test -v -cover ./...




.PHONY: up test