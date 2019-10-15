DEFAULT_TARGET: build

VERSION := $(shell git describe --tags --abbrev=0 --exact-match $(git rev-parse HEAD) 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS := -X github.com/makkes/gitlab-cli/config.Version=$(VERSION)

PLATFORMS := windows linux darwin
os = $(word 1, $@)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=build
BINARY_NAME=gitlab
SRCS := $(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -v -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gitlab

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p $(BUILD_DIR)
	GOOS=$(os) GOARCH=amd64 go build -v -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(os)_amd64 ./cmd/gitlab

.PHONY: install
install:
	go install ./cmd/gitlab

.PHONY: lint
lint:
	golangci-lint run
.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./$(BUILD_DIR)

.PHONY: release
release: windows linux darwin
