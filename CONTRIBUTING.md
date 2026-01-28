# Contributing to SemaMesh

First off, thank you for considering contributing to SemaMesh! It's people like you that make the AI ecosystem safer for everyone.

## 🛠 Development Environment
To work on SemaMesh, you will need:
* **Go 1.23+**
* **Docker** (for BPF compilation)
* **Kind or Minikube** (for local Kubernetes testing)
* **Clang/LLVM**

## 🤝 How to Contribute
1. **Fork the Repo:** Create your own branch from `main`.
2. **Setup BPF:** Ensure you compile the eBPF objects using `make build-bpf`.
3. **Tests:** Ensure your Go code passes `go test ./...`.
4. **Pull Request:** Submit a PR with a clear description of the change.

## 🏗 Middleware Architecture
SemaMesh is designed to be extensible. If you want to add a new feature (like PII masking or caching), please add it as a **Middleware** in `internal/proxy/middleware.go` rather than modifying the core proxy logic.

## 📝 Commit Messages
We follow [Conventional Commits](https://www.conventionalcommits.org/):
* `feat:` for new features.
* `fix:` for bug fixes.
* `docs:` for documentation changes.