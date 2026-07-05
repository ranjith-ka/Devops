# Gateway API Quick Reference

Quick lookup guide for common Gateway API operations.

## Installation

```bash
# Install Gateway API CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# Install Envoy Gateway controller
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo update
helm install envoy-gateway envoy-gateway/gateway-helm \
  --namespace envoy-gateway-system \
  --create-namespace
```

## Quick Commands

### View Resources
```bash
# List all Gateways
kubectl get gateways -A

# List all HTTPRoutes
kubectl get httproutes -A

# List all GatewayClasses
kubectl get gatewayclasses

# Get detailed Gateway info
kubectl describe gateway <name> -n <namespace>

# Get HTTPRoute status
kubectl describe httproute <name> -n <namespace>
```

### Debugging
```bash
# Check Gateway controller logs
kubectl logs -n envoy-gateway-system -l app=envoy-gateway --tail=50

# Watch Gateway resource changes
kubectl get gateway -n <namespace> -w

# Export resource as YAML
kubectl get httproute <name> -n <namespace> -o yaml

# Check resource events
kubectl describe httproute <name> -n <namespace>
```

### Port Forwarding
```bash
# Forward local port to Gateway service
kubectl port-forward -n <namespace> svc/envoy-gateway-<gateway-name> 8080:80

# Test with curl
curl -H "Host: <hostname>" http://localhost:8080/path
```

### Testing Routes
```bash
# Simple GET request
curl -H "Host: demo.localhost" http://localhost:8080/

# Request with headers
curl -H "Host: demo.localhost" -H "X-Custom: value" http://localhost:8080/

# Follow redirects
curl -L -H "Host: demo.localhost" http://localhost:8080/

# HTTPS (insecure for self-signed certs)
curl -k -H "Host: demo.localhost" https://localhost:8443/

# Verbose output
curl -v -H "Host: demo.localhost" http://localhost:8080/
```

## Basic Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: my-gateway
spec:
  gatewayClassName: envoy
  listeners:
  - name: http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same
```

## Basic HTTPRoute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-route
spec:
  parentRefs:
  - name: my-gateway
    kind: Gateway
  hostnames:
  - "example.com"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: my-service
      port: 8080
```

## TLS Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: tls-gateway
spec:
  gatewayClassName: envoy
  listeners:
  - name: https
    port: 443
    protocol: HTTPS
    tls:
      mode: Terminate
      certificateRefs:
      - name: my-cert-secret
        kind: Secret
```

## Path-Based Routing

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: path-route
spec:
  parentRefs:
  - name: my-gateway
    kind: Gateway
  hostnames:
  - "example.com"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: api-service
      port: 8080
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: web-service
      port: 8080
```

## Header-Based Routing

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: header-route
spec:
  parentRefs:
  - name: my-gateway
    kind: Gateway
  rules:
  - matches:
    - headers:
      - name: X-Version
        value: v2
    backendRefs:
    - name: app-v2
      port: 8080
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: app-v1
      port: 8080
```

## Weighted Load Balancing

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: weighted-route
spec:
  parentRefs:
  - name: my-gateway
    kind: Gateway
  rules:
  - backendRefs:
    - name: app-v1
      port: 8080
      weight: 70
    - name: app-v2
      port: 8080
      weight: 30
```

## Cross-Namespace Routing

```bash
# 1. Create ReferencePolicy in backend namespace
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferencePolicy
metadata:
  name: allow-route
  namespace: backend-namespace
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: gateway-namespace
  to:
  - group: ""
    kind: Service
EOF

# 2. Reference service in HTTPRoute
# backendRefs:
# - name: backend-service
#   namespace: backend-namespace
#   port: 8080
```

## Create TLS Secret

```bash
# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/CN=example.com"

# Create Kubernetes secret
kubectl create secret tls my-cert --cert=cert.pem --key=key.pem -n my-namespace

# Verify
kubectl get secret -n my-namespace
```

## Common Issues & Fixes

| Issue | Solution |
|-------|----------|
| Gateway not getting IP | Check controller is running: `kubectl get pods -n envoy-gateway-system` |
| HTTPRoute not routing | Verify backend service exists: `kubectl get svc -n <namespace>` |
| Connection refused | Check port-forward is running: `ps aux \| grep port-forward` |
| TLS handshake error | Use `-k` flag with curl for self-signed certs |
| Cross-namespace not working | Create ReferencePolicy in backend namespace |

## Useful Links

- [Gateway API Docs](https://gateway.api.k8s.io/)
- [Envoy Gateway Docs](https://gateway.envoyproxy.io/)
- [API Reference](https://gateway.api.k8s.io/v1alpha2/api-types/gateway/)
- [HTTPRoute Reference](https://gateway.api.k8s.io/v1alpha2/api-types/httproute/)

## Cleanup

```bash
# Delete resources
kubectl delete namespace demo-apps
kubectl delete namespace api-backend

# Uninstall Envoy Gateway
helm uninstall envoy-gateway -n envoy-gateway-system
kubectl delete namespace envoy-gateway-system

# Delete CRDs (optional)
kubectl delete crd -l gateway.networking.k8s.io/
```

---

See README.md for detailed tutorial with step-by-step examples.
