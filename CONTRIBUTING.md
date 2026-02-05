# Contributing to SemaMesh 🚀

First off, thank you for considering contributing to SemaMesh! It's people like you who make the AI infrastructure safer for everyone.

SemaMesh is an open-source project that adheres to the CNCF philosophy. We welcome contributions of all forms: code, documentation, bug reports, feature requests and even creative hacking solutions 🦹🏼.

# 📜 Code of Conduct
By participating in this project, you agree to abide by the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). Please read it to understand the expected behavior.

----

# 🛠️ Development Environment
SemaMesh is a hybrid project involving Go (**User Space**) and C (**eBPF Kernel Space**).

**Prerequisites**

* Go: v1.25+
* Docker: For building the hybrid image.
* Kubernetes Cluster: Kind, Minikube, or a remote dev cluster.
* Clang/LLVM: For compiling eBPF object files (if working on /bpf).

**Quick Start**

1. Fork & Clone
   ```
   git clone https://github.com/YOUR_USERNAME/SemaMesh.git
   cd SemaMesh
   ```
2. Build the Project
   We use Docker to handle the complexity of compiling both Go and eBPF.
   ```agsl
   docker build -t semamesh:v0.5.4 .
   ```
3. Run Locally (Kind)
   To test your changes in a real environment:
   ```
   # Load image into Kind
   kind load docker-image semamesh:v0.5.4 --name semamesh-lab
   
   # Deploy
   kubectl apply -f deploy/install.yaml
   ```

-----
# 🗺️ Project Structure
To help you navigate the codebase:

| Directory      | Description                                                    |
|----------------|----------------------------------------------------------------|
| `cmd/ `        | Main entry points (the binary).                                |
| `pkg/identity` | **Layer 8 Logic**. K8s API watchers & IP-to-Namespace mapping. |
| `pkg/proxy`    | The HTTP proxy that intercepts traffic.                        |
| `pkg/audit`    | Structured logging and compliance handling.                    |
| `pkg/metrics`  | Prometheus metric definitions and counters.                    |
| `pkg/sniffer`  | **Deep Packet Inspection**. OpenAI protocol parsers & logic.   |
| `bpf/`         | **Kernel Space**. C code for eBPF redirection.                 |
| `deploy/`      | Kubernetes manifests (install.yaml).                           |
| `dashboards/`  | Grafana JSON models for visualization.                         |
----

# 🏗️ Architecture & Middleware
SemaMesh is designed to be extensible.

**Adding New Features**:
If you want to add a feature like **PII Masking**, **Rate Limiting**, please implement it as **Middleware**.

* **Location**: `pkg/proxy/` or `pkg/sniffer/`
* **Guideline**: Do not modify the core `main.go` loop unless necessary. Chain your middleware in `pkg/proxy/handler.go` to keep the architecture clean.
----

# ✅ Pull Request Process
1. Sign Your Work (DCO): We require all commits to be signed off to comply with the [Developer Certificate of Origin](https://developercertificate.org).
   ```
   git commit -s -m "feat: add new pii masking middleware"
   ```
2. Atomic Commits: Keep PRs focused on a single issue.
3. Tests: New features must include unit tests.
4. Review: All PRs require approval from at least one maintainer.
----

# 📝 Commit Messages
We follow [Conventional Commits](https://www.conventionalcommits.org/):
* `feat`: A new feature
* `fix`: A bug fix
* `docs`: Documentation only changes
* `style`: Formatting, missing semi-colons, etc.
* `refactor`: A code change that neither fixes a bug nor adds a feature
* `test`: Adding missing tests
* `chore:` Build process or auxiliary tool changes
----

# 🐛 Reporting Bugs
Open a [GitHub Issue](https://github.com/semamesh/SemaMesh/issues) with:
* Your Kubernetes version (kubectl version).
* SemaMesh logs (kubectl logs -l app=semamesh).
* Steps to reproduce.
* Happy Hacking! 💻

----
