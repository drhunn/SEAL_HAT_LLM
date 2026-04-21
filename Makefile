APP_NAME := harness
DB_DSN ?= postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable

.PHONY: build run verify verify-soft verify-strict bootstrap-db check fmt tidy test pycheck pytest

build:
	go build ./...

run:
	go run ./cmd/harness -config config/runtime.example.toml

bootstrap-db:
	bash ./scripts/bootstrap_verify_db.sh '$(DB_DSN)'

verify: verify-soft

verify-soft:
	go run ./cmd/verify -config config/runtime.example.toml -mode soft

verify-strict:
	go run ./cmd/verify -config config/runtime.example.toml -mode strict

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
