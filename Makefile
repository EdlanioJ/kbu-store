.PHONY: test start build swag mock init-migration migrate-up migrate-down env

DATABASE="postgresql://postgres:root@db:5432/kbu_store?sslmode=disable"

test:
	go test ./... -cover

start:
	go run ./app/main.go

build:
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o kbu-store ./app/main.go

swag:
	swag init -g "./app/main.go" -d "./" -o "./app/docs"

mock:
	mockery --output "./domain/mocks" --dir "./" --all

init-migration:
	migrate create -ext sql -dir db/migration -seq init_schema

migrate-up:
	migrate -path db/migration -database ${DATABASE} -verbose up

migrate-down:
	migrate -path db/migration -database ${DATABASE} -verbose down

env:
	cp .env.example app.env