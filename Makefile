# Variables
PKG_NAME := viperenv

# Commands
GO := go
GOFMT := gofmt
GOLINT := golint
GOVET := go vet
GOTEST := go test
BINARY := $(PKG_NAME)

# Flags
GOFLAGS := -v
GOFMTFLAGS := -s
GOLINTFLAGS := -set_exit_status
GOVETFLAGS := -all
GOTESTFLAGS := -v

.PHONY: all
all: build test

.PHONY: build
build:
	$(GO) build $(GOFLAGS) -o $(BINARY)

.PHONY: test
test:
	$(GOTEST) $(GOTESTFLAGS) ./...

.PHONY: fmt
fmt:
	$(GOFMT) $(GOFMTFLAGS) -w .

.PHONY: lint
lint:
	$(GOLINT) $(GOLINTFLAGS) ./...

.PHONY: vet
vet:
	$(GOVET) $(GOVETFLAGS) ./...

.PHONY: clean
clean:
	rm -f $(BINARY)

	