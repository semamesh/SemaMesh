# --- Stage 1: Builder ---
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependency files first (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Waypoint Proxy (The Brain)
# CGO_ENABLED=0 creates a static binary (runs on any Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/waypoint ./cmd/waypoint/main.go

# Build the Loader (The Agent that loads eBPF)
# Note: We assume you have a loader.go. If not, we can run the proxy directly
# but usually, you need a small setup tool. For this MVP, we'll assume the
# Proxy *is* the entrypoint and loads eBPF internally if you merged them.
# If separate, build both. Let's assume a single binary for simplicity here.

# --- Stage 2: Runtime ---
# We use a distroless image for security (no shell, no unused apps)
# Or alpine if you want to debug. Let's use Alpine for the MVP.
FROM alpine:3.19

WORKDIR /root/

# Install necessary system libraries for networking (optional but good for debug)
RUN apk add --no-cache iproute2 ca-certificates

# Copy the Go Binary
COPY --from=builder /bin/waypoint .

# Copy the Compiled eBPF Object File
# CRITICAL: We built this in the previous step using Docker!
COPY bpf/sema_redirect.o /usr/lib/bpf/sema_redirect.o

# Expose the Waypoint Port
EXPOSE 15001

# Run the Proxy
CMD ["./waypoint"]