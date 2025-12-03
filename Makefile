PROTO_FILES=$(wildcard proto/service_*.proto)
SERVICES=$(patsubst proto/service_%.proto,%,$(PROTO_FILES))
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
