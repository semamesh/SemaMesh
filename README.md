![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)
![eBPF](https://img.shields.io/badge/Data_Plane-eBPF-ff69b4)
![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.25+-326CE5?style=flat&logo=kubernetes)

# 🕸SemaMesh ☁️🛡️
1. [ ] The "Layer 8" Service Mesh for AI Governance.
2. [ ] Observability, Identity, and Control for the AI Era.
----

## 📖 Introduction
**SemaMesh** (short for **Semantic Mesh**) is a specialized proxy designed to solve the "Shadow AI" problem in Kubernetes.

While traditional meshes (Istio, Linkerd) route packets at Layer 4 (TCP) or Layer 7 (HTTP), **SemaMesh** operates at **Layer 8 (Context)**. It sees Tokens, Costs, and Intent.

It sits between your microservices and external LLM providers (OpenAI, Anthropic, Gemini), providing deep observability, identity-aware cost tracking, and audit compliance without complex sidecars.

## 🚀 Why SemaMesh?
Feature | The Problem (Without SemaMesh)                                     | The Solution (With SemaMesh)                                                                  
--- |--------------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
💰 **Financial FinOps** | "Why is our OpenAI bill $5,000 this month?"                        | **Namespace-Level** Cost Attribution. See exactly which team/service is spending money.       | 
🕵️ **Data Security** | Developers might accidentally send PII (emails, names) to public LLMs. | **PII Redaction & Audit**. Automatically scrub sensitive data before it leaves the cluster.   | 
🆔 **Identity** | LLM providers only see one API Key.                                | **K8s Identity Mapping**. SemaMesh resolves the source Pod IP to a Namespace/Service Account. | 
⚖️ **Governance** | No record of what was asked.                                       | **The "Black Box.**" A structured audit log of every prompt and completion for compliance.    |
----

## 🏗️ Architecture
SemaMesh runs as a **Deployment** inside your Kubernetes cluster.
1. **The Identity Watcher**: Connects to the Kubernetes API to build a real-time map of Pod IPs to Namespaces.
2. **The Semantic Proxy**: Intercepts HTTP traffic to LLM providers, parses the JSON payload, and calculates token usage.
3. **The Auditor**: Asynchronously logs structured JSON events to dis (`/var/log/semamesh/audit.log`)
4. **The Reporter**: Exposes Prometheus metrics (`semamesh_llm_cost_est_total`) for Grafana visualization.
5. **eBPF Core (Preview)**: A compiled eBPF datapath is included for future transparent redirection (Coming in v0.6.0).

![SemaMesh Architecture](images/SemaMesh-Visio.png)

--------- 

## ⚡ Quick Start Guide
**Prerequisites**
* Kubernetes Cluster (Kind, Minikube, EKS, GKE)
* `kubectl` installed
* An OpenAI API Key

### 1. Installation
Deploy SemaMesh into your K8s cluster using the provided manifest.

```
# 1. Clone the repository
git clone https://github.com/semamesh/SemaMesh.git
cd SemaMesh

# 2. Build the image (if running locally)
docker build -t semamesh:v0.5.4 .
kind load docker-image semamesh:v0.5.4 --name semamesh-lab

# 3. Deploy (⚠️ Edit deploy/install.yaml to set your OpenAI API Key first!)
kubectl apply -f deploy/install.yaml
```

### 2. Verify Deployment
Ensure the proxy is running and has connected to the Kubernetes API.
```
kubectl logs -l app=semamesh
# Output should show:
# ⚡ Connected to Kubernetes API. Watching Pods...
# 🛡️ Identity Awareness: KUBERNETES (Real)
# 📊 Starting Metrics Server on :9090/metrics
```

### 3. Send a Test Request
You can use any pod inside the cluster to test the proxy.
```
# Launch a temporary curl pod
kubectl run test-client --rm -i --tty --image=curlimages/curl -- sh

# Send a request through SemaMesh
curl http://semamesh.default.svc.cluster.local:80/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-YOUR-'OpenAI'API-KEY" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Hello World!"}]
  }'
```

## 📊 Observability
**Prometheus Metrics**
SemaMesh exposes rich metrics at `:9090/metrics`.

Metric Name | Description | Labels                   
--- | --- |--------------------------|
`semamesh_llm_tokens_total` | Count of tokens (prompt vs completion). | `namespace, model, type` | 
`semamesh_llm_cost_est_total` | Estimated cost in USD based on public pricing. | `namespace, model`       |
`semamesh_http_requests_total` | Volume of requests and HTTP status codes. | `namespace, status`      | 

**Audit Logs**

SemaMesh writes a structured `NDJSON` audit log to `/var/log/semamesh/audit.log`. Even failed requests (401/429) are logged (Just in case, from my tests)
```
{
  "timestamp": "2026-02-04T20:46:52Z",
  "namespace": "finance-service",
  "model": "gpt-4",
  "prompt_text": "Analyze this transaction...",
  "completion_text": "The transaction appears valid.",
  "total_tokens": 150,
  "cost_usd": 0.03
}
```

### 💡 Dashboards: 
Import the pre-built dashboard from `/dashboards/semamesh-overview.json` into your Grafana instance to visualize real-time AI spend and token usage.

----

## 🛣️ Roadmap
* v0.5.0 (Alpha): ✅ Explicit Proxy, K8s Identity, Cost Metrics, Structured Audit Logs.
* v0.6.0 (Beta): 🚧 eBPF Transparent Redirection. (Eliminates the need to change application URLs).
* v0.7.0: Policy Engine (Rate limiting per namespace, blocking PII).
* v0.8.0: Intelligent Routing (Semantic Caching, Model Distillation) & Economics.
* v1.0.0: Production Release (Zero Trust Agents 🤞🏻).
----

## 🤝 Contributing
We welcome contributions! Please see [CONTRIBUTING.md](https://github.com/semamesh/SemaMesh/blob/2f349e2a04bc70b67f247dc91ba2d8727f721273/CONTRIBUTING.md) for details on how to set up your development environment.

## 📄 License
SemaMesh is open-source software licensed under the Apache 2.0 License.
