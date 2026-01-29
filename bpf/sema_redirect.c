#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <sys/socket.h>
#include <linux/in.h>

// --- STRUCT DEFINITIONS (MUST BE AT THE TOP) ---

// Key to identify a specific connection
struct sock_key {
    __u32 sip4;   // Source IP
    __u32 dip4;   // Destination IP
    __u32 sport;  // Source Port
    __u32 dport;  // Destination Port
    __u32 family; // Protocol Family (TCP/UDP)
};

// Value to store the original destination
struct origin_info {
    __u32 ip;     // Original IP
    __u32 port;   // Original Port
};

// --- MAP DEFINITIONS ---

// Helper macro to define maps (Standard libbpf syntax)
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 65535);
    __type(key, struct sock_key);
    __type(value, struct origin_info);
} proxy_map SEC(".maps");

// --- MAIN PROGRAM ---

SEC("sk_msg")
int sema_redirect(struct sk_msg_md *msg)
{
    // Now the compiler knows what 'struct sock_key' looks like!
    struct sock_key key = {0};

    // Extract key details from the message metadata
    key.sip4 = msg->remote_ip4;
    key.dip4 = msg->local_ip4;
    key.sport = msg->remote_port;
    key.dport = msg->local_port;
    key.family = msg->family;

    // Look up if this connection should be redirected
    struct origin_info *val = bpf_map_lookup_elem(&proxy_map, &key);

    if (val) {
        // Redirect the traffic to the proxy port
        long ret = bpf_msg_redirect_hash(msg, &proxy_map, &key, BPF_F_INGRESS);
        if (ret != 0) {
            bpf_printk("SemaMesh: Redirect failed: %ld\n", ret);
        }
        return (int)ret;
    }

    return SK_PASS;
}

char _license[] SEC("license") = "GPL";