API_DOCS = docs/api

KEEPER_VERSION ?= 0.1.0
SERVER_VERSION ?= 0.1.0

BUILD_DATE ?= $(shell date +%F\ %H:%M:%S)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)

PROTO_SRC = proto/keeper/v1
PROTO_FILES = auth notification secrets
PROTO_DST = pkg/$(PROTO_SRC)

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
.PHONY: keeper

# swag init -g ./internal/httpserver/backend.go --output docs/api
server: ## build server
	rm -rf $(API_DOCS)
	
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

clean: ## remove build artifacts
	rm -rf cmd/keeper/keeper cmd/server/server cmd/staticlint/staticlint
.PHONY: clean

unit-tests: ## run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
.PHONY: unit-tests

# godoc: ### show public packages documentation using godoc
# 	@echo "Project documentation is available at:"
# 	@echo "http://127.0.0.1:3000/pkg/github.com/ex0rcist/metflix/pkg/\n"
# 	@godoc -http=:3000 -play
# .PHONY: godoc

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
