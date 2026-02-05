# Security Policy

The SemaMesh community takes security seriously. We appreciate your efforts to responsibly disclose your findings, and will make every effort to acknowledge your contributions.

## Supported Versions

As a CNCF Sandbox project, we currently prioritize security updates for the latest available release.

| Version | Supported          |
|---------| ------------------ |
| v0.5.x  | :white_check_mark: |
| < v0.5  | :x:                |

## Reporting a Vulnerability

If you believe you have found a security vulnerability in SemaMesh, **please do not open a public issue.**

Instead, please send a report via email to:
**semamesh009@gmail.com** (Subject: `[SECURITY] SemaMesh Vulnerability Report`)

Please include as much of the following as possible:
* A description of the vulnerability.
* Steps to reproduce the issue (including policy YAMLs or specific LLM prompts used).
* Potential impact (e.g., Privilege Escalation, Policy Bypass, Denial of Service).
* Any proof-of-concept code or screenshots.

### Response Timeline
* **Acknowledgment:** We will aim to respond to your report within 48 hours.
* **Assessment:** We will determine if the finding is valid and assess its severity.
* **Fix:** We will prepare a patch and release it as a priority update.
* **Disclosure:** Once the fix is released, we will publicly announce the vulnerability and credit you for the discovery (unless you prefer to remain anonymous).

## Threat Model & Scope

SemaMesh operates with high privileges (`CAP_SYS_ADMIN`, `CAP_NET_ADMIN`) to load eBPF programs and manage CRIU checkpoints.

**In Scope:**
* **Audit Evasion**: Methods to send traffic through the proxy without it being recorded in the Audit Log.
* **Identity Spoofing**: Tricking the Identity Watcher into attributing cost/usage to the wrong Kubernetes Namespace.* **Privilege Escalation:** Utilizing the SemaMesh DaemonSet to gain unauthorized root access to the host node.
* **Privilege Escalation**: Utilizing the container capabilities to gain unauthorized access to the host node.
* **Denial of Service:** Crashing the Waypoint Proxy or Metrics Server via malformed traffic.

**Out of Scope:**
* Vulnerabilities in the underlying Kubernetes cluster (e.g., compromised kubelet) unless caused by SemaMesh.
* Attacks requiring physical access to the node.
* Spam or social engineering attacks.
* "Policy Bypass" (Note: Blocking policies are not yet enforced in v0.5.x).