# 🗺️ SemaMesh Roadmap

SemaMesh is evolving from a Layer 4/7 Traffic Monitor into a **Layer 8 (Semantic) Operating System** for Autonomous Agents. Our roadmap focuses on three pillars: **Governance**, **Economics**, and **Safety**.

##  v0.4.0: The initial Foundation
- [ ] ~~**eBPF Datapath:** Transparent interception of pod egress traffic using `sock_ops` and `sk_msg`.~~
- [x] **Waypoint Proxy:** High-performance Go-based sidecarless proxy for HTTP/JSON parsing.
- [x] **Intent Blocking:** Regex and Keyword-based blocking for "Destructive Intent" (e.g., `DELETE`, `DROP TABLE`).
- [x] **Stateful Pause (Alpha):** Integration with CRIU to freeze pods upon policy violation.
- [x] **Mock LLM Testing:** End-to-End smoke testing pipeline for CI/CD verification.

---

## ✅ 🔭 v0.5.0: Deep Observability (Current Status)
*Focus: Bringing "Layer 8" metrics to Kubernetes.*

- [x] **Prometheus Exporter**: Export semamesh_llm_tokens_total, semamesh_llm_cost_est_total, and traffic health metrics.
- [x] **Structured Audit Logging**: Log every prompt/response pair to /var/log/semamesh/audit.log with redacted PII (JSON format).
- [x] **Grafana Dashboard Pack**: Pre-built dashboards for "Cost per Namespace" and "Traffic Volume."
- [x] **Identity Mapping**: Resolving Pod IPs to Kubernetes Namespaces for attribution.

## 🚧 v0.6.0: The Cognitive Firewall (Q3 2026)
*Focus: Making the mesh "invisible" via eBPF.*

- [ ] **Transparent Redirection (eBPF):** 
  - Activate the `sock_ops` and `sk_msg` BPF programs.
  - Eliminate the need for applications to change their URL. Redirect `api.openai.com` -> `localhost:8080` at the kernel level.
- [ ] **DaemonSet Architecture**: Refactor from Deployment to Node-based DaemonSet.

## 🛡️ v0.7.0: The Cognitive Firewall (Q3 2026)
Focus: Protecting data and preventing jailbreaks.

- [ ] **PII Redaction Middlewar**e: "Snap-in" filter to detect and redact Emails, SSNs, and API Keys before they leave the cluster.
- [ ] **Intent Blocking**: Regex and Keyword-based blocking for "Destructive Intent" (e.g., `DELETE`, `DROP TABLE`).
- [ ] **Stateful Pause (CRIU)**: Integration with CRIU to physically "freeze" a pod upon policy violation.

## 🔀 v0.8.0: Intelligent Routing & Economics (Q4 2026)
Focus: Reducing AI costs via "Model Distillation Routing."
- [ ] **Semantic Caching**: If a user asks a question that is semantically similar (Vector Cosine Similarity) to a previous question, return the cached answer without hitting the paid LLM API.
- [ ] **Complexity-Based Routing**: Use a local classifier (BERT-tiny) to judge prompt complexity.
  - Simple Prompt -> Route to local `Ollama`/`Llama-3` (Free).
  - Complex Prompt -> Route to `GPT-4` (Paid).
---

## 🆔 v1.0.0: Zero Trust Agents (Hopefully Q4 2026)
Focus: Production Hardening.
- [ ] **Agent Identity Attestation** (**SPIFFE**): Bind Kubernetes Service Accounts to SPIFFE IDs.
- [ ] **mTLS for LLM Gateways**: Enforce mutual TLS between the Agent and the Waypoint Proxy.
- [ ] **Multi-Cluster Federation**: Enforce global AI safety policies across hybrid cloud (AWS/GCP/On-Prem).