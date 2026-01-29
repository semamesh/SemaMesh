#!/bin/bash
set -e

# --- 1. CONFIGURATION & PATHS ---
# Automatically find the project root relative to this script
REPO_ROOT="$( cd -- "$(dirname "$0")/.." >/dev/null 2>&1 ; pwd -P )"

# Define manifest paths relative to root
KIND_CONFIG="$REPO_ROOT/hack/kind-config.yaml"
MOCK_LLM="$REPO_ROOT/test/mock-llm.yaml"
RBAC_YAML="$REPO_ROOT/deploy/rbac.yaml"
CRD_DIR="$REPO_ROOT/config/crd/bases"
DAEMONSET_YAML="$REPO_ROOT/deploy/daemonset.yaml"
POLICY_EXAMPLE="$REPO_ROOT/examples/sample-policy.yaml"

echo "-------------------------------------------------------"
echo "🚀 SemaMesh Pre-Submission Smoke Test"
echo "Project Root: $REPO_ROOT"
echo "-------------------------------------------------------"

# --- 2. CLUSTER CHECK ---
if ! kind get clusters | grep -q "semamesh-lab"; then
    echo "🏗️ Cluster 'semamesh-lab' not found. Creating it..."
    kind create cluster --config "$KIND_CONFIG"
else
    echo "✅ Using existing 'semamesh-lab' cluster."
fi

# ADD THIS LINE HERE:
echo "🎯 Switching context to kind-semamesh-lab..."
kubectl config use-context kind-semamesh-lab

echo "🚚 Loading semamesh-agent:latest into Kind..."
kind load docker-image semamesh-agent:latest --name semamesh-lab

# --- 3. DEPLOYMENT PHASE ---
echo "🏗️ Creating Namespace..."
# We use 'apply' here so it doesn't error out if it already exists
kubectl create namespace semamesh-system --dry-run=client -o yaml | kubectl apply -f -

echo "📦 Applying SemaMesh Infrastructure..."
kubectl apply -f "$RBAC_YAML"
kubectl apply -f "$CRD_DIR"
kubectl apply -f "$DAEMONSET_YAML"

echo "🤖 Deploying Mock LLM Server..."
kubectl apply -f "$MOCK_LLM"

echo "⏳ Waiting for SemaMesh Agent Pods to be created..."
# Adding a small sleep to ensure K8s has started the pod before we look for it
sleep 5
kubectl rollout status daemonset/semamesh-node-agent -n semamesh-system --timeout=90s

# --- 4. EXECUTION PHASE ---
echo "🔍 Running Traffic Interception Tests..."

# Find the agent pod
AGENT_POD=$(kubectl get pod -n semamesh-system -l app=semamesh -o jsonpath="{.items[0].metadata.name}")

echo "[TEST A] Sending 'Safe' prompt via $AGENT_POD..."
# Exec into the pod and run the curl command
kubectl exec -n semamesh-system "$AGENT_POD" -- curl -s -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello world"}' | grep "SUCCESS"

echo "✅ Safe request passed."

echo "[TEST B] Applying Semantic Policy..."
# Ensure the policy file exists in your repo
kubectl apply -f "$POLICY_EXAMPLE"
sleep 2 # Give the agent a second to see the new policy

echo "🚫 [TEST C] Simulating Destructive Intent..."
# This should return the 'Paused' message we saw earlier
BLOCK_RESPONSE=$(kubectl exec -n semamesh-system "$AGENT_POD" -- curl -s -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"prompt": "I want to delete the production database"}')

if [[ "$BLOCK_RESPONSE" == *"SemaMesh Policy Violation"* ]]; then
    echo "✅ Success: Destructive intent blocked by SemaMesh!"
else
    echo "❌ Error: Policy was not enforced."
    exit 1
fi