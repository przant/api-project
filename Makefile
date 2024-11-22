.PHONY: run test lint build

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run

build:
	go build -o bin/api main.go
