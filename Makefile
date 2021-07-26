.PHONY: test start.http start.grpc start.kafka build swag mock migrate.create migrate.up migrate.down gen env

DATABASE="postgresql://postgres:root@db:5432/kbu_store?sslmode=disable"

test:
	go test ./...

test.cover:
	go test ./... -v -covermode=count -coverprofile=coverage.out

start.http:
	go run ./app/main.go http

start.grpc:
	go run ./app/main.go grpc

start.kafka:
	go run ./app/main.go kafka

build:
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o kbu-store ./app/main.go

swag:
	swag init -g "./app/infrastructure/http/server.go" -d "./" -o "./app/infrastructure/http/docs"

mock:
	mockery --output "./app/utils/mocks" --dir "./app/domain/" --all
	mockery --output "./app/utils/mocks" --dir "./app/interfaces/" --all

migrate.create:
	migrate create -ext sql -dir app/db/migration -seq init_schema

migrate.up:
	migrate -path app/db/migration -database ${DATABASE} -verbose up

migrate.down:
	migrate -path app/db/migration -database ${DATABASE} -verbose down

gen:
	protoc --proto_path=app/infrastructure/grpc app/infrastructure/grpc/protofiles/*.proto --go_out=. --go-grpc_out=.

env:
	cp .env.example app.env
