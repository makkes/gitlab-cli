DEFAULT_TARGET: build

VERSION := next
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/makkes/gitlab-cli/config.Version=$(VERSION) -X github.com/makkes/gitlab-cli/config.Commit=$(COMMIT)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=build
PKG_DIR=/gitlab-cli
BINARY_NAME=gitlab-cli
SRCS := $(shell find . -type f -name '*.go' -not -path './vendor/*')

$(BUILD_DIR)/$(BINARY_NAME): $(SRCS)
	mkdir -p $(BUILD_DIR)
	go build -v -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: install
install:
	go install -v -ldflags '$(LDFLAGS)'

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./$(BUILD_DIR)
