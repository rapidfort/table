
.PHONY: all build test clean deps run example

all: deps build test

deps:
	go mod tidy
	go mod download

build:
	go build ./...

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	go clean
	rm -f coverage.out coverage.html

run:
	go run example/main.go