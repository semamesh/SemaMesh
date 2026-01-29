# --- Stage 1: Build eBPF Artifacts (The "Linux Lab") ---
FROM ubuntu:22.04 AS bpf-builder

# Install tools to compile C code for BPF
# We include build-essential and libbpf-dev to support the compilation
RUN apt-get update && apt-get install -y \
    clang llvm libbpf-dev \
    linux-headers-generic \
    make \
    build-essential

WORKDIR /build

# FIX 1: Create a symlink for the 'asm' headers so Clang can find them on ARM64
# This maps /usr/include/aarch64-linux-gnu/asm -> /usr/include/asm
RUN ln -s /usr/include/$(uname -m)-linux-gnu/asm /usr/include/asm

# Copy only the C source code
COPY bpf/ /build/bpf/

# FIX 2: Compile the eBPF program with explicit include paths
# We add -I flags to tell Clang where to find standard headers (sys/socket.h, etc.)
RUN clang -O2 -g -target bpf \
    -I/usr/include/$(uname -m)-linux-gnu \
    -c bpf/sema_redirect.c \
    -o bpf/sema_redirect.o

# --- Stage 2: Build the Go Proxy (The "Muscle") ---
# FIX 3: Use Go 1.25 (or latest) to match your local go.mod requirement
FROM golang:1.25 AS go-builder

WORKDIR /app

# Copy Go dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Waypoint binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/waypoint ./cmd/waypoint/main.go

# --- Stage 3: Final Runtime Image ---
FROM alpine:3.19

WORKDIR /root/

# ✅ ADD basic networking tools for debugging (ip, etc.) & curl here so it's always available for testing/healthchecks
RUN apk add --no-cache iproute2 ca-certificates curl

# 1. Copy the Go binary from Stage 2
COPY --from=go-builder /bin/waypoint .

# 2. Copy the Compiled eBPF object from Stage 1
COPY --from=bpf-builder /build/bpf/sema_redirect.o /usr/lib/bpf/sema_redirect.o

# (Optional) Expose the port your agent listens on.
# Adjust this number if your code uses something else (e.g., 15000, 8000)
EXPOSE 8080

# Start the agent
CMD ["./waypoint"]