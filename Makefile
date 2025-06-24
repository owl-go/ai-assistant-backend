# Makefile for AI Assistant Backend

# 变量定义
BINARY_NAME=ai-assistant-backend
BUILD_DIR=build
MAIN_FILE=main.go

# Go相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# 构建参数
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# 默认目标
.PHONY: all
all: clean build

# 构建应用
.PHONY: build
build:
	@echo "构建应用..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行应用
.PHONY: run
run:
	@echo "运行应用..."
	$(GORUN) $(MAIN_FILE)

# 测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "更新依赖..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "代码检查..."
	$(GOCMD) vet ./...

# 代码质量检查
.PHONY: lint
lint: fmt vet
	@echo "代码质量检查完成"

# 数据库初始化
.PHONY: init-db
init-db:
	@echo "初始化数据库..."
	@chmod +x init_db.sh
	./init_db.sh

# 开发模式运行
.PHONY: dev
dev:
	@echo "开发模式运行..."
	@chmod +x start.sh
	./start.sh

# 安装应用
.PHONY: install
install: build
	@echo "安装应用..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "应用已安装到 /usr/local/bin/$(BINARY_NAME)"

# 卸载应用
.PHONY: uninstall
uninstall:
	@echo "卸载应用..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "应用已卸载"

# 生成API文档
.PHONY: docs
docs:
	@echo "生成API文档..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_FILE); \
		echo "API文档已生成"; \
	else \
		echo "请先安装 swag: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(BINARY_NAME):latest .

# 运行Docker容器
.PHONY: docker-run
docker-run:
	@echo "运行Docker容器..."
	docker run -p 8080:8080 --env-file config.env $(BINARY_NAME):latest

# 帮助信息
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  build        - 构建应用"
	@echo "  run          - 运行应用"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试并生成覆盖率报告"
	@echo "  clean        - 清理构建文件"
	@echo "  deps         - 安装依赖"
	@echo "  deps-update  - 更新依赖"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo "  lint         - 代码质量检查"
	@echo "  init-db      - 初始化数据库"
	@echo "  dev          - 开发模式运行"
	@echo "  install      - 安装应用"
	@echo "  uninstall    - 卸载应用"
	@echo "  docs         - 生成API文档"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"
	@echo "  help         - 显示帮助信息" 