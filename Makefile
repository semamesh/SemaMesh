# Variables
IMAGE_NAME ?= semamesh:latest
BIN_DIR = bin
BPF_DIR = bpf

.PHONY: all build-bpf build-go build-controller build-waypoint docker-build deploy clean help

all: build-bpf build-go ## Build everything (BPF + Go binaries)

## --- Build Targets ---

build-bpf: ## Compile eBPF C code using Docker for cross-platform compatibility
	@echo "==> Compiling eBPF bytecode..."
	docker run --rm -v $(shell pwd):/src -w /src \
		quay.io/cilium/ebpf-builder:1.6.0 \
		clang -O2 -g -target bpf -c $(BPF_DIR)/sema_redirect.c -o $(BPF_DIR)/sema_redirect.o

build-go: build-controller build-waypoint ## Build all Go binaries

build-controller: ## Build the main K8s Controller
	@echo "==> Building Sema Controller..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/sema-controller ./cmd/main.go

build-waypoint: ## Build the Waypoint Proxy
	@echo "==> Building Waypoint Proxy..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/waypoint ./cmd/waypoint/main.go

## --- Docker & Deployment ---

docker-build: build-bpf ## Build the Docker image for the cluster
	@echo "==> Building Docker image $(IMAGE_NAME)..."
	docker build -t $(IMAGE_NAME) .

deploy: ## Apply CRDs and DaemonSet to the current K8s context
	@echo "==> Deploying to Kubernetes..."
	kubectl apply -f config/crd/bases/
	kubectl apply -f deploy/daemonset.yaml

clean: ## Remove binaries and compiled objects
	@echo "==> Cleaning up..."
	rm -rf $(BIN_DIR)
	rm -f $(BPF_DIR)/*.o

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'