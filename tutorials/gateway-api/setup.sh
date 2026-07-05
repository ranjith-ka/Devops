#!/bin/bash
# Quick Setup Script for Gateway API Tutorial
# Run this script to quickly set up the entire Gateway API environment

set -e

echo "🚀 Starting Gateway API Setup..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Create Kind Cluster
echo -e "${BLUE}[1/5] Creating Kind cluster...${NC}"
kind delete cluster --name k8s 2>/dev/null || true
cd /Users/ranjith.a/code/github/ranjith/Devops
kind create cluster --config kind/config.yaml --name k8s
echo -e "${GREEN}✓ Kind cluster created${NC}\n"

# Step 2: Install Gateway API CRDs
echo -e "${BLUE}[2/5] Installing Gateway API CRDs...${NC}"
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml
sleep 5
echo -e "${GREEN}✓ Gateway API CRDs installed${NC}\n"

# Step 3: Install Envoy Gateway
echo -e "${BLUE}[3/5] Installing Envoy Gateway controller...${NC}"
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo update
helm install envoy-gateway envoy-gateway/gateway-helm \
  --namespace envoy-gateway-system \
  --create-namespace
echo -e "${GREEN}✓ Envoy Gateway installed${NC}\n"

# Step 4: Wait for controller
echo -e "${BLUE}[4/5] Waiting for Envoy Gateway controller to be ready...${NC}"
kubectl wait --timeout=5m -n envoy-gateway-system \
  --for=condition=Progressing=True \
  deployment/envoy-gateway 2>/dev/null || echo "Warning: timeout waiting for deployment"
sleep 10
echo -e "${GREEN}✓ Envoy Gateway is ready${NC}\n"

# Step 5: Verify installation
echo -e "${BLUE}[5/5] Verifying installation...${NC}"
echo -e "${YELLOW}Cluster Info:${NC}"
kubectl cluster-info
echo -e "\n${YELLOW}Nodes:${NC}"
kubectl get nodes
echo -e "\n${YELLOW}Gateway Classes:${NC}"
kubectl get gatewayclasses
echo -e "\n${YELLOW}Envoy Gateway Pods:${NC}"
kubectl get pods -n envoy-gateway-system

echo -e "\n${GREEN}✅ Setup Complete!${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Read the full tutorial: tutorials/gateway-api/README.md"
echo "2. Follow Part 2 to deploy your first application"
echo "3. Create a Gateway and HTTPRoute to test routing"

