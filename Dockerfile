# ==========================================
# Stage 1: Build eBPF Artifacts (The "Kernel Lab")
# ==========================================
FROM ubuntu:22.04 AS bpf-builder

# Install tools to compile C code for BPF
RUN apt-get update && apt-get install -y \
    clang llvm libbpf-dev \
    linux-headers-generic \
    make \
    build-essential

WORKDIR /build

# FIX: Create a symlink for 'asm' headers so Clang can find them on ARM64/x86
RUN ln -s /usr/include/$(uname -m)-linux-gnu/asm /usr/include/asm

# Copy the BPF source code
COPY bpf/ /build/bpf/

# Compile the eBPF program
# We add -I flags to tell Clang where to find standard headers
RUN clang -O2 -g -target bpf \
    -I/usr/include/$(uname -m)-linux-gnu \
    -c bpf/sema_redirect.c \
    -o bpf/sema_redirect.o

# ==========================================
# Stage 2: Build the Go Proxy (The "User Space")
# ==========================================
FROM golang:1.25 AS go-builder

WORKDIR /app

# Copy Go dependencies first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the SemaMesh binary (Statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/semamesh ./cmd/semamesh/main.go

# ==========================================
# Stage 3: Runtime Image (The "Production Artifact")
# ==========================================
FROM alpine:3.19

WORKDIR /root/

# Install tools
RUN apk add --no-cache iproute2 ca-certificates curl

# 1. Create the Audit Log Directory explicitly 📂 <--- ADD THIS LINE
RUN mkdir -p /var/log/semamesh

# 2. Copy the Go binary
COPY --from=go-builder /bin/semamesh .

# 3. Copy the Compiled eBPF object
COPY --from=bpf-builder /build/bpf/sema_redirect.o /usr/lib/bpf/sema_redirect.o

EXPOSE 8080 9090

CMD ["./semamesh", "--dev=false"]