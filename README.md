# 🕸️ SemaMesh: The Semantic AI Service Mesh
**SemaMesh** is a next-generation, sidecarless service mesh built specifically for the governance, security, and financial oversight of Autonomous AI Agents in Kubernetes.

By moving beyond simple Layer 4/7 networking into Layer 8+ Semantic Networking, SemaMesh understands not just where a packet is going, but the intent of the AI agent sending it.

## ✨ Key Features
- **Sidecarless eBPF Interception**: Transparently hijacks AI traffic at the kernel level. No sidecars, no manual pod injection, and near-zero latency overhead.
- **Stateful Pause (The Kill-Switch)**: Uses CRIU to freeze high-risk agents mid-thought. If an agent tries to delete a namespace or bypass safety protocols, its RAM is snapshotted to disk for human review.
- **Semantic Quota Management**: Define token budgets via CRDs. Prevent recursive agent loops from draining your OpenAI/Anthropic wallet.
- **Modular Middleware**: A Go-based Waypoint Proxy that allows you to "snap in" new features like PII redaction or LLM response caching.

## 🏗️ Architecture
SemaMesh consists of three primary components that work in harmony:

1. **The Brain (Control Plane)**: A Kubernetes Controller written in Go that manages SemaPolicy and SemaTokenQuota CRDs.
2. **The Trap (eBPF Interceptor)**: A C program loaded into the Linux Kernel that redirects outbound LLM traffic to the local Waypoint.
3. **The Muscle (Waypoint Proxy)**: A high-performance Go proxy that analyzes "Reasoning Traces," counts tokens, and enforces security policies.

# 🚀 Getting Started
## 1. Project Structure

```
semamesh/
├── api/v1alpha1/       # Custom Resource Definitions (Go Types)
├── bpf/                # eBPF C-code for traffic redirection
├── cmd/
│   ├── main.go         # Control Plane (Manager)
│   └── waypoint/       # Data Plane (Proxy)
├── internal/
│   ├── controller/     # Reconciler logic for CRDs
│   └── proxy/          # Middleware chain (TokenGuards, Policies)
│   └── agent/          # eBPF loader and manager
└── deploy/             # K8s Manifests (DaemonSets, RBAC)
```
## 2. Compilation

To compile the eBPF kernel program on a non-Linux machine (like macOS), use our Docker-based builder:

```
docker run --rm -v $(pwd):/src -w /src \
quay.io/cilium/ebpf-builder:1.6.0 \
clang -O2 -g -target bpf -c bpf/sema_redirect.c -o bpf/sema_redirect.o
```

## 3. Deployment

```
# Apply the "Bank" and "Law" CRDs
kubectl apply -f config/crd/bases/

# Deploy the SemaMesh Node Agent (DaemonSet)
kubectl apply -f deploy/daemonset.yaml```
```

## Security Warning
- **Important**: SemaMesh requires `CAP_SYS_ADMIN` to load eBPF programs and use CRIU for checkpointing. Ensure you understand the security implications before deploying in production.
- In our `daemonset.yaml`, we used `privileged: true`, that because the mesh uses eBPF & it requires these permissions. For production, consider using more granular capabilities.


## 📜 Example Policy
Define a "Hard Freeze" for any agent that attempts to destroy infrastructure or exceeds a $50 token budget:
```
YAML
apiVersion: semamesh.io/v1alpha1
kind: SemaPolicy
metadata:
name: infrastructure-safety-gate
spec:
rules:
- name: "prevent-unauthorized-deletion"
intentMatches: ["delete namespace", "terminate node"]
riskLevel: "Critical"
action: "PAUSE"  # Triggers CRIU Freeze
pauseSettings:
timeout: "30m"
notify: "https://hooks.slack.com/services/T123/B456"
```

## 🛠️ Technical Deep Dive
**The "Stateful Pause" Flow**

When the Waypoint Proxy detects a Critical violation in a prompt:

1. It annotates the Pod: semamesh.io/action: PAUSE. 
2. The Sema Controller sees the annotation. 
3.The Controller calls the Kubelet Checkpoint API. 
3. The Agent process is frozen; its memory is saved to /var/lib/kubelet/checkpoints. 
4. A DevOps Architect uses kubectl sema approve to thaw the pod or terminate it.

## 🤝 Contributing
SemaMesh is a CNCF-style project. We follow the Middleware Pattern for extensibility. To add a new feature, simply add a new filter to the internal/proxy/middleware.go chain.