APP_NAME := harness

.PHONY: build run verify fmt tidy

build:
	go build ./...

run:
	go run ./cmd/harness -config config/runtime.example.toml

verify:
	go run ./cmd/verify -config config/runtime.example.toml

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy
