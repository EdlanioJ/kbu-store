.PHONY: test start-http start-grpc build swag mock create-migration migrate-up migrate-down grpc-gen env

DATABASE="postgresql://postgres:root@db:5432/kbu_store?sslmode=disable"

test:
	go test ./... -cover

start-http:
	go run ./main.go http

start-grpc:
	go run ./main.go grpc

build:
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o kbu-store ./main.go

swag:
	swag init -g "./application/http/server.go" -d "./" -o "./application/http/docs"

mock:
	mockery --output "./domain/mocks" --dir "./domain/" --all

create-migration:
	migrate create -ext sql -dir infra/db/migration -seq init_schema

migrate-up:
	migrate -path infra/db/migration -database ${DATABASE} -verbose up

migrate-down:
	migrate -path infra/db/migration -database ${DATABASE} -verbose down

grpc-gen:
	protoc --proto_path=application/grpc application/grpc/protofiles/*.proto --go_out=. --go-grpc_out=.

env:
	cp .env.example app.env