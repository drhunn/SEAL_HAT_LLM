APP_NAME := harness

.PHONY: build run verify check fmt tidy test pycheck

build:
	go build ./...

run:
	go run ./cmd/harness -config config/runtime.example.toml

verify:
	go run ./cmd/verify -config config/runtime.example.toml

check:
	go test ./...
	python -m compileall python

test:
	go test ./...

pycheck:
	python -m compileall python

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy
