# Gateway API Example Manifests

This directory contains ready-to-use example manifests for the Gateway API tutorial.

## Files

- `01-basic-gateway.yaml` - Simple Gateway and HTTPRoute example
- `02-advanced-routing.yaml` - Multi-path and header-based routing
- `03-tls-gateway.yaml` - Gateway with HTTPS/TLS support
- `04-cross-namespace.yaml` - Cross-namespace routing with ReferencePolicy
- `05-traffic-management.yaml` - Weighted load balancing
- `06-sample-apps.yaml` - Sample applications for testing

## Quick Start

Deploy the example namespace and apps:
```bash
kubectl create namespace demo-apps
kubectl apply -f 06-sample-apps.yaml
kubectl apply -f 01-basic-gateway.yaml

# Test
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &
curl -H "Host: demo.localhost" http://localhost:8080/
pkill -f "port-forward"
```

## Progressive Learning Path

1. **Start simple** - Use `01-basic-gateway.yaml` for basic HTTP routing
2. **Advance routing** - Use `02-advanced-routing.yaml` for more complex rules
3. **Add HTTPS** - Use `03-tls-gateway.yaml` for TLS termination
4. **Multi-namespace** - Use `04-cross-namespace.yaml` to learn ReferencePolicy
5. **Traffic control** - Use `05-traffic-management.yaml` for weighted routing

See README.md for detailed walkthrough instructions.
