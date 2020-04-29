PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NewApp \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=asd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=ascli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		go install  $(BUILD_FLAGS) ./cmd/pxbcli
		go install  $(BUILD_FLAGS) ./cmd/pxbd
		go install  $(BUILD_FLAGS) ./cmd/pxbrelayer

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Build
build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbcli.exe ./cmd/pxbcli
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbd.exe ./cmd/pxbd
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbrelayer.exe ./cmd/pxbrelayer
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbcli ./cmd/pxbcli
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbd ./cmd/pxbd
	go build -mod=readonly $(BUILD_FLAGS) -o build/pxbrelayer ./cmd/pxbrelayer
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify


# Local environment

build-docker-pxbdnode:
	$(MAKE) -C networks/local

localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/pxbd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/pxbd:Z lcnem/proximax-pegzone testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi

	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down
