# 设置变量
APP_NAME := imapi
BIN_DIR := ./bin
PROTO_DIR := ./proto

API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: all clean run api

all: clean proto build run

# 编译项目
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) main.go
	chmod +x $(BIN_DIR)/$(APP_NAME)

# 运行项目
run:
	@echo "Running $(APP_NAME)..."
	$(BIN_DIR)/$(APP_NAME)

# 清理构建文件
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# generate api proto
api:
	@echo "Compiling Protobuf files..."
	protoc --proto_path=./api \
		   --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
		   --go-grpc_out=paths=source_relative:./api \
		   --grpc-gateway_out=paths=source_relative:./api \
	       $(API_PROTO_FILES)
