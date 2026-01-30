# 📦 Release v0.1.0: "The Big Bang" 🌌
We are excited to announce the initial alpha release of SemaMesh, the first sidecarless service mesh designed specifically for the era of autonomous AI agents. This release lays the foundation for "Semantic Networking"—where the mesh understands the intent of the traffic it carries.

## 🚀 Key Highlights

### ⚡ Sidecarless eBPF Interception
We’ve eliminated the sidecar tax. By utilizing eBPF cgroup/connect4 hooks, SemaMesh transparently redirects AI traffic directly to the local Waypoint Proxy at the kernel level.
- Benefit: Lower latency and simplified operations. No more sidecar.istio.io/inject.

## ⏸️ Stateful Pause (CRIU)
The standout feature of this release. SemaMesh can now "freeze" a misbehaving agent mid-reasoning.
- Action: When a security policy is violated, the mesh triggers a CRIU (Checkpoint/Restore in Userspace) snapshot.
- Why: This allows human-in-the-loop verification before a recursive agent spends your budget or deletes your cluster.

## 💰 Semantic Token Quota
Built-in support for LLM Token Governance.
- Feature: Real-time token counting via the Waypoint Proxy.
- Control: Automatically block or rate-limit requests that exceed the semantic budget defined in SemaTokenQuota CRDs.

## 🛠️ Technical Specifications
- Core: Go 1.23+
- Kernel: eBPF (C) with CO-RE support.
- Kubernetes: Tested on v1.28 - v1.31.
- Platform: Currently optimized for Linux (Required for eBPF and CRIU).

## 🏗️ What’s in the box?
- `sema-controller`: The Kubernetes-native brain managing your policies.
- `sema-waypoint`: The high-performance L8 proxy with a modular middleware chain.
- `sema-interceptor`: The eBPF-powered "trap" that catches AI traffic.
CRDs: `SemaPolicy` and `SemaTokenQuota`.

## ⚠️ Alpha Notice
This is an experimental alpha release. It is intended for DevOps architects and AI researchers to explore the possibilities of semantic networking.
* "Do not run this in your high-frequency trading prod yet—unless you're feeling adventurous."

## 🤝 Contributors
A huge shoutout to everyone involved in the initial architecture and design of the semantic "Kill-Switch."

---

# Release v0.2.0: "Minor updates" 🛠️
Updated Adopters, Contributing,Governance & Maintainers guides.

---

# Release v0.3.0: "Bug fixes & smoke-test plan" 🐛
Updated the Data Plane, Control Plane, Intent Blocking, Kubectl components, Mock LLM Testing & added a smoke-test plan for new users to validate their installation.

---

# Release v0.4.0: "Roadmap" 🗺️
Added ROADMAP.md to track future updates and releases.