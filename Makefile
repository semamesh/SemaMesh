.PHONY: build-bpf build-proxy deploy

build-bpf:
	docker run --rm -v $(shell pwd):/src -w /src \
		quay.io/cilium/ebpf-builder:1.6.0 \
		clang -O2 -g -target bpf -c bpf/sema_redirect.c -o bpf/sema_redirect.o

build-proxy:
	go build -o bin/waypoint ./cmd/waypoint/main.go

deploy:
	kubectl apply -f config/crd/bases/
	kubectl apply -f deploy/daemonset.yaml