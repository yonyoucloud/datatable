.DEFAULT_GOAL:=help

.EXPORT_ALL_VARIABLES:

ifndef VERBOSE
.SILENT:
endif

SHELL=/bin/bash -o pipefail -o errexit

REPO_INFO ?= $(shell git config --get remote.origin.url)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
COMMIT_TIME ?= $(shell git show --pretty="%ci %cr" | head -1 | awk '{print $$1"/"$$2"/"$$3}')
BUILD_TIME ?= $(shell date "+%Y-%m-%d/%H:%M:%S/%z")
REGISTRY ?= github.com/yonyoucloud
VERSION ?= v1.0.1

datatable_LDFLAGS=-w -X github.com/yonyoucloud/datatable/pkg/version.Pkg=$(REPO_INFO) -X github.com/yonyoucloud/datatable/pkg/version.Version=$(VERSION) -X github.com/yonyoucloud/datatable/pkg/version.GitCommitSha=$(COMMIT_SHA) -X github.com/yonyoucloud/datatable/pkg/version.GitCommitTime=$(COMMIT_TIME) -X github.com/yonyoucloud/datatable/pkg/version.BuildTime=$(BUILD_TIME)

HOST_ARCH = $(shell which go >/dev/null 2>&1 && go env GOARCH)
HOST_OS = $(shell which go >/dev/null 2>&1 && go env GOOS)
ARCH ?= $(HOST_ARCH)
OS ?= $(HOST_OS)
# 打arm包时强制一下OS，否则Mac下生成的无法在linux arm中运行
OS = linux
ifeq ($(ARCH),)
    $(error mandatory variable ARCH is empty, either set it when calling the command or make sure 'go env GOARCH' works)
endif

help:  ## 请查看下面帮助
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

IMAGE ?= $(REGISTRY)/datatable:$(VERSION)-$(COMMIT_SHA)-$(ARCH)

.PHONY: image
image: clean-image ## 构建镜像
	echo "Building docker image ($(ARCH))..."
	@docker build \
		--no-cache \
		--build-arg TARGETARCH="$(ARCH)" \
		--build-arg COMMIT_SHA="$(COMMIT_SHA)" \
		-t $(IMAGE) deploy

.PHONY: push
push:  ## 推送镜像
	@docker push $(IMAGE)

.PHONY: clean-image
clean-image: ## 删除本地镜像
	echo "removing old image $(REGISTRY)/datatable:*"
	@docker rmi `docker images | grep "$(REGISTRY)/datatable" | awk '{print $$3}'` || true

.PHONY: build-be
build-be:  ## 编译 datatable 后端
	@cd backend && \
		rm -rf ./deploy/$(ARCH)/datatable && \
		go mod tidy && \
		CGO_ENABLED=0 \
		GOOS=$(OS) \
		GOARCH=$(ARCH) \
		go build -ldflags "$(datatable_LDFLAGS)" -o ../deploy/$(ARCH)/datatable ./cmd

.PHONY: build-fe
build-fe:  ## 编译 datatable 前端
	@npm run build

.PHONY: run-be
run-be:  ## 运行 datatable 后端
	@cd backend && \
		go mod tidy && \
		go run -ldflags "$(datatable_LDFLAGS)" ./cmd

.PHONY: run-fe
run-fe:  ## 运行 datatable 前端
	@npm run dev

.PHONY: deploy
deploy:  ## 部署 datatable
	@kubectl apply -f ./deploy/datatable.yaml