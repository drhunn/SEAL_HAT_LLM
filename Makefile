APP_NAME := harness

.PHONY: build run check fmt tidy test

build:
	go build ./...

run:
	go run ./cmd/harness -config config/runtime.example.toml

check:
	go test ./...

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy
