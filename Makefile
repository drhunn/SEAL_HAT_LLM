APP_NAME := harness

.PHONY: build run verify check fmt tidy test pycheck pytest

build:
	go build ./...

run:
	go run ./cmd/harness -config config/runtime.example.toml

verify:
	go run ./cmd/verify -config config/runtime.example.toml

check:
	go test ./...
	python -m compileall python
	python -m unittest discover -s python/tests

test:
	go test ./...
	python -m unittest discover -s python/tests

pycheck:
	python -m compileall python

pytest:
	python -m unittest discover -s python/tests

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy
