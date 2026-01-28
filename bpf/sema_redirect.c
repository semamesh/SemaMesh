#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <sys/socket.h>

// 1. Define a Map to store the "Original Destination"
// When we hijack the packet, we need to remember where it *wanted* to go,
// so our Proxy can look it up later.
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 65535);
    __type(key, struct sock_key);
    __type(value, struct origin_info);
} original_dst_map SEC(".maps");

// 2. The Hook: "cgroup/connect4"
// This triggers every time a container tries to start a TCP connection.
SEC("cgroup/connect4")
int sema_connect4(struct bpf_sock_addr *ctx) {
    // A. Filter: Ignore localhost traffic (don't redirect our own proxy!)
    if (ctx->user_ip4 == bpf_htonl(0x7F000001)) {
        return 1; // Allow
    }

    // B. Filter: Only redirect standard AI ports (e.g., HTTPS 443)
    // In a real mesh, we would use a map of "Allowed IPs" here.
    if (ctx->user_port != bpf_htons(443) && ctx->user_port != bpf_htons(80)) {
        return 1; // Allow non-web traffic to pass normally
    }

    // C. Save the Original Destination
    // We create a "Key" based on the source/dest pair
    struct sock_key key = {
        .sip = ctx->msg_src_ip4,
        .dip = ctx->user_ip4,
        .sport = 0, // In connect4, source port isn't set yet, so we rely on cookie/pid
        .dport = ctx->user_port
    };

    struct origin_info value = {
        .original_ip = ctx->user_ip4,
        .original_port = ctx->user_port
    };

    bpf_map_update_elem(&original_dst_map, &key, &value, BPF_ANY);

    // D. The Redirect (The Magic)
    // We rewrite the destination to Localhost (127.0.0.1)
    ctx->user_ip4 = bpf_htonl(0x7F000001);

    // We rewrite the destination port to our Waypoint Proxy (15001)
    ctx->user_port = bpf_htons(15001);

    // Print a debug message to the kernel trace pipe (view with: cat /sys/kernel/debug/tracing/trace_pipe)
    bpf_printk("SemaMesh: Redirected traffic to Waypoint Proxy :15001\n");

    return 1; // Allow the (now modified) connection to proceed
}

char _license[] SEC("license") = "GPL";