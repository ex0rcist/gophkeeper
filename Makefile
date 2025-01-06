API_DOCS = docs/api

KEEPER_VERSION ?= 1.0.0
SERVER_VERSION ?= 1.0.0

BUILD_DATE ?= $(shell date +%d.%m.%y)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_DIR = build

PROTO_SRC = proto/keeper/grpcapi
PROTO_FILES = users notification secrets health
PROTO_DST = pkg/$(PROTO_SRC)

PLATFORMS = \
    linux/amd64 \
    windows/amd64 \
    darwin/amd64 \
    darwin/arm64

help: ## display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

build: keeper server
.PHONY: build

keeper: ## build keeper
	go build \
		-ldflags "\
			-X 'main.buildVersion=$(KEEPER_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go

	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT=$(BUILD_DIR)/keeper-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build \
			-ldflags "\
				-X 'main.buildVersion=$(KEEPER_VERSION)' \
				-X 'main.buildDate=$(BUILD_DATE)' \
				-X 'main.buildCommit=$(BUILD_COMMIT)' \
			" \
			-o $$OUTPUT \
			cmd/keeper/*.go || exit 1; \
	done

.PHONY: keeper

server: ## build server	
	go build \
		-ldflags "\
			-X 'main.buildVersion=$(SERVER_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go
.PHONY: server

staticlint: ## build static lint
	go build -o cmd/$@/$@ cmd/$@/*.go
.PHONY: staticlint

certs: ## generate certs
	cd cert && ./gen.sh > /dev/null
.PHONY: certs

clean: ## remove build artifacts
	rm -rf cmd/keeper/keeper cmd/server/server cmd/staticlint/staticlint
.PHONY: clean

unit-tests: ## run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
.PHONY: unit-tests

proto: $(PROTO_FILES) ## generate gRPC protobuf bindings
.PHONY: proto

$(PROTO_FILES): %: $(PROTO_DST)/%

$(PROTO_DST)/%:
	protoc \
		--proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_DST) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DST) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/$(notdir $@).proto
