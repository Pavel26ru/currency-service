PROTO_SRC=proto/currency/currency_service.proto
PROTO_OUT=.

PROTOC_GEN_GO=protoc-gen-go
PROTOC_GEN_GO_GRPC=protoc-gen-go-grpc

GEN_TEST_DATA_SCRIPT=./currency/internal/scripts/generate_test_data.go

.PHONY: all build run test proto

# air для hot reload https://github.com/air-verse/air

all: build

# Альтернатива таргету proto (+ generate.go)
gen:
	go generate ./...

# Лучше убрать latest из версия внешних зависимостей.
install-tools:
	@echo "Checking and installing necessary tools..."
	@if ! [ -x "$$(command -v $(PROTOC_GEN_GO))" ]; then \
		echo "Installing protoc-gen-go..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \ # для учебного проекта норм latest
	fi
	@if ! [ -x "$$(command -v $(PROTOC_GEN_GO_GRPC))" ]; then \
		echo "Installing protoc-gen-go-grpc..."; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	fi

proto: install-tools
	@echo "Generating gRPC and Protobuf code..."
	protoc --proto_path=proto --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)

build: proto
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./... -cover

generate-test-data:
	@echo "Generating and inserting test data into the database..."
	go run $(GEN_TEST_DATA_SCRIPT)

mocks:
	go generate ./...