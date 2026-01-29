# Contributing to SemaMesh

First off, thank you for considering contributing to SemaMesh! It's people like you who make the AI infrastructure safer for everyone.

## Getting Started
1. **Fork** the repository on GitHub.
2. **Clone** your fork locally.
3. Install dependencies: `go`, `clang`, `llvm`, and `docker`.
4. Run `make all` to ensure the project builds on your machine.

## Technical Standards
- **DCO**: All commits must be signed off (`git commit -s`) to comply with the Developer Certificate of Origin.
- **Testing**: Ensure your eBPF and Go changes include unit tests.
- **Style**: Follow standard `gofmt` for Go and Linux kernel style for C/eBPF code.

## How to Contribute
- **Bug Reports**: Open a GitHub Issue with logs and reproduction steps.
- **Feature Requests**: Open an Issue to discuss the design before writing code.
- **Pull Requests**: Submit PRs against the `main` branch. All PRs require at least one maintainer approval.

## 🏗 Middleware Architecture
SemaMesh is designed to be extensible. If you want to add a new feature (like PII masking or caching), please add it as a **Middleware** in `internal/proxy/middleware.go` rather than modifying the core proxy logic.

## 📝 Commit Messages
We follow [Conventional Commits](https://www.conventionalcommits.org/):
* `feat:` for new features.
* `fix:` for bug fixes.
* `docs:` for documentation changes.