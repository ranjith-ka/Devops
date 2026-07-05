# Gateway API Tutorial Index

Welcome to the comprehensive Gateway API tutorial for Kubernetes! This extended tutorial covers everything from basic setup to advanced traffic management.

## 📚 Tutorial Structure

### Main Documentation
- **[README.md](./README.md)** - Complete step-by-step tutorial (7 parts, ~22KB)
  - Part 1: Setup & Installation
  - Part 2: Basic HTTP Routing
  - Part 3: Advanced Routing
  - Part 4: TLS/HTTPS
  - Part 5: Cross-Namespace Routing
  - Part 6: Traffic Management
  - Part 7: Troubleshooting

### Quick Reference
- **[QUICK_REFERENCE.md](./QUICK_REFERENCE.md)** - Command cheatsheet and common patterns
  - Installation commands
  - Debugging commands
  - YAML templates
  - Common issues & fixes

### Example Manifests
All example manifests are ready-to-use and can be applied directly:

1. **[06-sample-apps.yaml](./06-sample-apps.yaml)** - Start here!
   - Deploy sample applications (http-app, api-app)
   - Creates demo-apps namespace
   - Essential for testing

2. **[01-basic-gateway.yaml](./01-basic-gateway.yaml)** - Basics
   - Simple Gateway listening on port 80
   - Basic HTTPRoute for hostname routing
   - Best for learning the core concepts

3. **[02-advanced-routing.yaml](./02-advanced-routing.yaml)** - Intermediate
   - Path-based routing (/api vs /)
   - Header-based routing (X-Version matching)
   - Multiple hostname routing

4. **[03-tls-gateway.yaml](./03-tls-gateway.yaml)** - Security
   - HTTPS/TLS termination
   - Certificate management
   - HTTP to HTTPS redirect

5. **[04-cross-namespace.yaml](./04-cross-namespace.yaml)** - Advanced
   - Cross-namespace service routing
   - ReferencePolicy for security
   - Multi-namespace architectures

6. **[05-traffic-management.yaml](./05-traffic-management.yaml)** - Advanced
   - Weighted load balancing (70/30 splits)
   - Canary deployments
   - User-based canary routing
   - Blue-green deployments

### Setup Script
- **[setup.sh](./setup.sh)** - Automated one-command setup
  - Creates kind cluster with Gateway API support
  - Installs Gateway API CRDs
  - Installs Envoy Gateway controller
  - Verifies everything is working

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended for beginners)
```bash
cd /Users/ranjith.a/code/github/ranjith/Devops/tutorials/gateway-api
chmod +x setup.sh
./setup.sh
```

Then follow [README.md Part 2](./README.md#part-2-basic-http-routing) for first application.

### Option 2: Manual Setup
```bash
# Navigate to Devops directory
cd /Users/ranjith.a/code/github/ranjith/Devops

# Create cluster
kind delete cluster --name k8s
kind create cluster --config kind/config.yaml --name k8s

# Install Gateway API
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# Install Envoy Gateway
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo update
helm install envoy-gateway envoy-gateway/gateway-helm \
  --namespace envoy-gateway-system \
  --create-namespace
```

### Option 3: Progressive Learning
Start simple and progress through examples:

```bash
# 1. Deploy sample apps
kubectl create namespace demo-apps
kubectl apply -f 06-sample-apps.yaml

# 2. Create basic gateway
kubectl apply -f 01-basic-gateway.yaml

# 3. Test basic routing
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &
curl -H "Host: demo.localhost" http://localhost:8080/
pkill -f "port-forward"

# 4. Add advanced features
kubectl apply -f 02-advanced-routing.yaml
# ... then 03, 04, 05 as you progress
```

## 📖 Learning Path

**Beginner** → **Intermediate** → **Advanced**

```
Week 1: Basics
├─ Read: README.md Part 1-2
├─ Do: Deploy sample apps (06-sample-apps.yaml)
└─ Do: Create basic gateway (01-basic-gateway.yaml)

Week 2: Routing
├─ Read: README.md Part 3
├─ Do: Apply advanced-routing.yaml
└─ Do: Test path and header-based routing

Week 3: Security & Multi-Namespace
├─ Read: README.md Part 4-5
├─ Do: Setup TLS (03-tls-gateway.yaml)
└─ Do: Cross-namespace routing (04-cross-namespace.yaml)

Week 4: Advanced Traffic Management
├─ Read: README.md Part 6
└─ Do: Weighted routing and canary (05-traffic-management.yaml)

Week 5: Troubleshooting & Production
├─ Read: README.md Part 7
└─ Practice: Common issues and debugging
```

## 🎯 What You'll Learn

✅ Gateway API concepts and architecture
✅ Installing and configuring Envoy Gateway
✅ Creating and managing Gateway resources
✅ HTTPRoute advanced routing patterns
✅ TLS/HTTPS termination
✅ Cross-namespace routing with ReferencePolicy
✅ Traffic management (weighted, canary, blue-green)
✅ Troubleshooting common issues
✅ Migrating from Ingress to Gateway API

## 📊 File Sizes & Content Overview

| File | Size | Purpose |
|------|------|---------|
| README.md | 22KB | Complete tutorial with all 7 parts |
| QUICK_REFERENCE.md | 6KB | Command cheatsheet |
| 06-sample-apps.yaml | 2KB | Test applications |
| 01-basic-gateway.yaml | 1KB | Basic example |
| 02-advanced-routing.yaml | 3KB | Advanced routing |
| 03-tls-gateway.yaml | 3KB | TLS/HTTPS |
| 04-cross-namespace.yaml | 3KB | Multi-namespace |
| 05-traffic-management.yaml | 5KB | Traffic management |
| setup.sh | 2KB | Automated setup |

**Total: ~47KB of documentation + 18KB of examples**

## 🔗 Prerequisites

- Docker (running)
- kind v0.20+
- kubectl v1.31+
- helm v3+
- OpenSSL (for TLS examples)

Verify:
```bash
docker --version
kind version
kubectl version --client
helm version
openssl version
```

## 🆘 Help & Troubleshooting

### Common Questions

**Q: Can I use this with an existing cluster?**
A: Yes! Skip the kind cluster creation and install Gateway API CRDs + Envoy Gateway directly.

**Q: Do I need to remove my Ingress controller?**
A: No, you can run both simultaneously for testing.

**Q: How do I switch back to Ingress?**
A: Just delete Gateway and HTTPRoute resources. See Cleanup section in README.md.

**Q: Can I use a different Gateway controller?**
A: Yes! Try Kong, Cilium, or Nginx instead of Envoy Gateway.

### Debugging Commands

```bash
# Check controller status
kubectl get pods -n envoy-gateway-system

# View controller logs
kubectl logs -n envoy-gateway-system -l app=envoy-gateway --tail=50

# Inspect resources
kubectl describe gateway <name> -n <namespace>
kubectl describe httproute <name> -n <namespace>

# See QUICK_REFERENCE.md for more commands
```

## 📚 Additional Resources

- [Official Gateway API Docs](https://gateway.api.k8s.io/)
- [Envoy Gateway Documentation](https://gateway.envoyproxy.io/)
- [Gateway API GitHub](https://github.com/kubernetes-sigs/gateway-api)
- [Envoy Proxy Docs](https://www.envoyproxy.io/docs/)

## 🎓 After This Tutorial

Once you've completed this tutorial, consider:

1. **Explore Other Route Types**
   - GRPCRoute for gRPC services
   - TCPRoute for TCP traffic
   - TLSRoute for TLS passthrough

2. **Advanced Features**
   - Custom filters
   - BackendPolicy for advanced routing
   - TCPRoute with arbitrary protocols

3. **Integration**
   - Set up Prometheus metrics
   - Configure Grafana dashboards
   - Implement Fluentd/ELK logging

4. **Production Deployment**
   - Multi-cluster setup
   - High availability configuration
   - Security policies and RBAC

5. **Automation**
   - GitOps with Flux or ArgoCD
   - Terraform/Pulumi for Gateway API
   - Custom controllers

## 📝 Notes

- All examples use the `envoy` GatewayClass. Adjust if using a different controller.
- Examples use self-signed certificates for TLS. Use proper certificates in production.
- Demo apps use 2 replicas. Scale as needed for your testing.
- All examples include inline comments with testing instructions.

## ✅ Verification Checklist

After completing the tutorial:
- [ ] Kind cluster is running with Gateway API support
- [ ] Gateway API CRDs are installed
- [ ] Envoy Gateway controller is running
- [ ] Sample apps are deployed and responding
- [ ] Basic HTTPRoute is working
- [ ] Advanced routing (paths, headers) works
- [ ] TLS/HTTPS is configured
- [ ] Cross-namespace routing works
- [ ] Weighted routing is functioning
- [ ] You can troubleshoot common issues

---

**Happy Learning!** 🚀

Start with `setup.sh` or follow [README.md](./README.md) for the complete step-by-step guide.
