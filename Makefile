# Base template, do not edit

# These variables can be overriden on the command line
# Example: make server APP_VERSION=1.2.3-rc1
export GOPRIVATE?=github.com/softsense
DOCKER_REGISTRY?=gcr.io/softsenseio
APP_VERSION?=$(shell git rev-parse --short HEAD)

# Use "make <target> VERBOSE=true" to run tests with verbose flag.
ifdef VERBOSE
	TEST_OPTS="-v"
endif

.PHONY: all
all: build lint test

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: test
test:
	go test $(TEST_OPTS) -race ./... -coverprofile=coverage.tmp && grep -vE 'cmd/.*' coverage.tmp > coverage.out
	go tool cover -html ./coverage.out -o coverage.html
	go tool cover -func ./coverage.out -o coverage.txt

.PHONY: integration-test
integration-test:
	go test $(TEST_OPTS) ./... -tags integration

.PHONY: db-test
db-test:
	go test $(TEST_OPTS) ./... -tags db

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: generate
	go build ./...

.PHONY: clean
clean:
	go clean ./...

# Repo-specific additions below
