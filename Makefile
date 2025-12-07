# Binary name
BINARY_NAME=uranus
VERSION?=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

PROTO_FILES=$(wildcard proto/service_*.proto)
SERVICES=$(patsubst proto/service_%.proto,%,$(PROTO_FILES))

# ==================== Build Commands ====================

.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build complete: ./$(BINARY_NAME)"

.PHONY: build-all
build-all: build-linux build-darwin build-windows
	@echo "✅ All builds complete!"

.PHONY: build-linux
build-linux:
	@echo "🐧 Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .

.PHONY: build-darwin
build-darwin:
	@echo "🍎 Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .

.PHONY: build-windows
build-windows:
	@echo "🪟 Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: install
install: build
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installed! Run '$(BINARY_NAME) --help' to get started."

.PHONY: uninstall
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled!"

.PHONY: clean
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "✅ Clean complete!"

# ==================== Development Commands ====================

.PHONY: run
run:
	go run . $(ARGS)

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...

# ==================== Proto Commands ====================

init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	#go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	#go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	#go install github.com/google/wire/cmd/wire@latest
	go mod tidy

api:
	@echo "found service: $(SERVICES)"
	@for service in $(SERVICES); do \
		echo "generating service: $$service"; \
		mkdir -p pb/$$service; \
		protoc --proto_path=./proto \
		       --proto_path=./grpc_third_party \
		       --go_out=./pb/$$service \
		       --go_opt=paths=source_relative \
		       --go-grpc_out=./pb/$$service \
		       --go-grpc_opt=paths=source_relative \
		       proto/service_$$service.proto; \
	done
	@echo "✓ finished generate API"

.PHONY: validate
validate:
	@echo "generate validation service: $(SERVICES)"
	@for service in $(SERVICES); do \
		echo "generating validation service: $$service"; \
		mkdir -p pb/$$service; \
		protoc --proto_path=./proto \
		       --proto_path=./grpc_third_party \
		       --go_out=./pb/$$service \
		       --go_opt=paths=source_relative \
		       --validate_out=lang=go:./pb/$$service \
		       --validate_opt=paths=source_relative \
		       proto/service_$$service.proto; \
	done
	@echo "✓ Completed generate validation"

generate:
	make api
	make validate

#sudo mv uranus /usr/local/bin/