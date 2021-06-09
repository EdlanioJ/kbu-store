
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
