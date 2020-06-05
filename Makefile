# Base template, do not edit

.PHONY: all
all: build lint test integration-test

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: test
test:
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: integration-test
integration-test:
	go test ./... -tags integration

.PHONY: db-test
db-test:
	go test ./... -tags db

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: generate
	go build ./...

# Repo-specific additions below