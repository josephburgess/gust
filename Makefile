.PHONY: build test test-unit test-integration test-cover lint clean all

BINARY_NAME=gust
GO=go
GOTEST=$(GO) test
GOLINT=golangci-lint

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/gust

test: test-unit test-integration

test-unit:
	$(GOTEST) -v ./... -short

test-integration:
	$(GOTEST) -v ./... -run Integration

test-cover:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

lint:
	$(GOLINT) run ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

all: clean lint test build
