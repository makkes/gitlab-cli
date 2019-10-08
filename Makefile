DEFAULT_TARGET: release

VERSION := 3.3.0-dev
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/makkes/gitlab-cli/config.Version=$(VERSION) -X github.com/makkes/gitlab-cli/config.Commit=$(COMMIT)

PLATFORMS := windows linux darwin
os = $(word 1, $@)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=build
PKG_DIR=/gitlab-cli_$(GOOS)
BINARY_NAME=gitlab
SRCS := $(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p $(BUILD_DIR)
	GOOS=$(os) GOARCH=amd64 go build -v -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(os)_amd64

.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: install
install:
	go install -v -ldflags '$(LDFLAGS)'

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
