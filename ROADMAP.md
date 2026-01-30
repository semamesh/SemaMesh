# 🗺️ SemaMesh Roadmap

SemaMesh is evolving from a Layer 4/7 Traffic Monitor into a **Layer 8 (Semantic) Operating System** for Autonomous Agents. Our roadmap focuses on three pillars: **Governance**, **Economics**, and **Safety**.

## ✅ v0.4.0: The Foundation (Current Status)
- [x] **eBPF Datapath:** Transparent interception of pod egress traffic using `sock_ops` and `sk_msg`.
- [x] **Waypoint Proxy:** High-performance Go-based sidecarless proxy for HTTP/JSON parsing.
- [x] **Intent Blocking:** Regex and Keyword-based blocking for "Destructive Intent" (e.g., `DELETE`, `DROP TABLE`).
- [x] **Stateful Pause (Alpha):** Integration with CRIU to freeze pods upon policy violation.
- [x] **Mock LLM Testing:** End-to-End smoke testing pipeline for CI/CD verification.

---

## 🔭 v0.5.0: Deep Observability (Q2 2026)
*Focus: Bringing "Layer 8" metrics to Kubernetes monitoring.*

- [ ] **Prometheus Exporter for Tokens:** - Export `llm_token_usage_total`, `llm_cost_est_hourly`, and `prompt_latency_seconds` metrics.
    - Differentiate between "Prompt Tokens" (Input) and "Completion Tokens" (Output) for cost tracking.
- [ ] **Structured Audit Logging:** - Log every prompt/response pair to an external sink (S3/Elasticsearch) with redacted PII, tied to the Kubernetes Service Account.
- [ ] **Grafana Dashboard Pack:** - Pre-built dashboards for "Cost per Namespace" and "Top Consumer Agents."

## 🛡️ v0.6.0: The Cognitive Firewall (Q3 2026)
*Focus: protecting data and preventing jailbreaks.*

- [ ] **PII Redaction Middleware:** - "Snap-in" WASM filter to detect and redact Emails, SSNs, and API Keys *before* they leave the cluster (using regex + local lightweight NLP).
- [ ] **Jailbreak Detection:** - Integration with tools like *Rebuff* or *Lakera* to detect "DAN" (Do Anything Now) or prompt injection attacks.
- [ ] **Presidio Integration:** - Native support for Microsoft Presidio for enterprise-grade data anonymization.

## 🔀 v0.7.0: Intelligent Routing & Economics (Q4 2026)
*Focus: Reducing AI costs via "Model Distillation Routing."*

- [ ] **Semantic Caching:** - If a user asks a question that is *semantically similar* (Vector Cosine Similarity) to a previous question, return the cached answer without hitting the paid LLM API.
- [ ] **Complexity-Based Routing:** - Use a small, local classifier (e.g., BERT-tiny) to judge prompt complexity.
    - **Simple Prompt** -> Route to local `Ollama/Llama-3` (Free).
    - **Complex Prompt** -> Route to `GPT-4` (Paid).

## 🆔 v0.8.0: Agent Identity (SPIFFE) (2027)
*Focus: Zero Trust for Agents.*

- [ ] **Agent Identity Attestation:** - Bind Kubernetes Service Accounts to SPIFFE IDs.
    - Ensure that *only* specific Agents can access specific Models (e.g., "Only the `Finance-Bot` can query the `GPT-4-Finance` model").
- [ ] **mTLS for LLM Gateways:** - Enforce mutual TLS between the Agent and the Waypoint Proxy.

---

## 🔮 Future Vision (v1.0+)
- **Interactive Forensics:** "Un-pause" a frozen agent in a sandbox environment to debug its behavior safely.
- **Multi-Cluster Federation:** Enforce global AI safety policies across hybrid cloud (AWS/GCP/On-Prem).