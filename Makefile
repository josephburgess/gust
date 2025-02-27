.PHONY: build test test-unit test-integration test-cover lint clean install all

BINARY_NAME=gust
GO=go
GOTEST=$(GO) test
GOLINT=golangci-lint
INSTALL_DIR=/usr/local/bin

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/gust

test: test-unit test-integration

test-v: test-unit-v test-integration-v

test-unit:
	$(GOTEST) ./... -short

test-integration:
	$(GOTEST) ./... -run Integration

test-unit-v:
	$(GOTEST) -v ./... -short

test-integration-v:
	$(GOTEST) -v ./... -run Integration

test-cover:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

lint:
	$(GOLINT) run ./...

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

all: clean lint test build
