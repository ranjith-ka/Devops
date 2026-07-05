# Gateway API Setup for Kind Cluster

This guide covers installing and configuring Gateway API on your local kind cluster.

## Prerequisites

- kind cluster running with K8s 1.31.6+
- kubectl installed
- helm installed

## Installation Steps

### 1. Recreate Kind Cluster with Gateway API Support

The `config.yaml` has been updated with Gateway API feature gates enabled. Recreate your cluster:

```bash
# Delete existing cluster if running
kind delete cluster --name k8s

# Create new cluster with Gateway API support
kind create cluster --config kind/config.yaml --name k8s
```

Verify the cluster is running:
```bash
kubectl cluster-info
kubectl get nodes
```

### 2. Install Gateway API CRDs

Install the Gateway API custom resource definitions:

```bash
# Install stable Gateway API CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# Verify CRDs are installed
kubectl get crds | grep gateway
```

Expected output:
```
gatewayclasses.gateway.networking.k8s.io
gateways.gateway.networking.k8s.io
httproutes.gateway.networking.k8s.io
grpcroutes.gateway.networking.k8s.io
tlsroutes.gateway.networking.k8s.io
tcproutes.gateway.networking.k8s.io
udproutes.gateway.networking.k8s.io
referencepolicies.gateway.networking.k8s.io
backendtlspolicies.gateway.networking.k8s.io
```

### 3. Install a Gateway Controller (Choose One)

You need a controller to implement the Gateway API. Popular options:

#### Option A: Envoy Gateway (Recommended)
```bash
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo update

helm install envoy-gateway envoy-gateway/gateway-helm \
  --namespace envoy-gateway-system \
  --create-namespace
```

#### Option B: Nginx Gateway (Alternative)
```bash
helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update

helm install nginx-gateway nginx-stable/nginx-gateway \
  --namespace nginx-gateway-system \
  --create-namespace
```

#### Option C: Cilium (if using Cilium for networking)
```bash
helm repo add cilium https://helm.cilium.io
helm repo update

helm install cilium cilium/cilium \
  --namespace kube-system \
  --set gatewayAPI.enabled=true
```

### 4. Verify Installation

Check if the Gateway controller is running:

```bash
# For Envoy Gateway
kubectl get pods -n envoy-gateway-system

# For Nginx Gateway
kubectl get pods -n nginx-gateway-system

# Check available GatewayClasses
kubectl get gatewayclasses
```

### 5. Create Your First Gateway

Here's an example Gateway and HTTPRoute to replace Ingress:

#### Create a Gateway (gateway-example.yaml)
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: my-gateway
  namespace: default
spec:
  gatewayClassName: envoy  # Use 'nginx' if using Nginx Gateway
  listeners:
    - name: http
      port: 80
      protocol: HTTP
    - name: https
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name: my-cert
            kind: Secret
```

#### Create an HTTPRoute (httproute-example.yaml)
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-app-route
  namespace: default
spec:
  parentRefs:
    - name: my-gateway
      kind: Gateway
  hostnames:
    - "awesome-http.example.com"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /dev
      backendRefs:
        - name: my-app
          port: 8080
```

Apply these:
```bash
kubectl apply -f gateway-example.yaml
kubectl apply -f httproute-example.yaml
```

### 6. Test the Gateway

```bash
# Get the Gateway IP/Port
kubectl get gateway my-gateway

# For port-forwarding (local testing)
kubectl port-forward svc/my-gateway 8080:80

# Test with curl
curl -H "Host: awesome-http.example.com" http://localhost:8080/dev
```

## Key Differences from Ingress to Gateway API

| Aspect | Ingress | Gateway API |
|--------|---------|------------|
| **API Resources** | Single Ingress | Gateway + HTTPRoute/TCPRoute/etc |
| **Port Mapping** | Limited | Full control via Listeners |
| **TLS** | Basic SNI support | Advanced TLS policies |
| **Routing** | Simple path/host rules | Complex matching (headers, params, etc) |
| **Authorization** | Via annotations | ReferencePolicy for fine-grained control |
| **Controller** | Nginx/Traefik/etc | Multiple options |

## Migration from Ingress to Gateway API

Quick comparison if you're migrating:

### Old Ingress:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
spec:
  rules:
    - host: awesome-http.example.com
      http:
        paths:
          - path: /dev
            backend:
              service:
                name: my-app
                port:
                  number: 8080
```

### New Gateway API:
```yaml
# Step 1: Create Gateway (one per cluster/namespace)
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

---
# Step 2: Create HTTPRoute (replaces Ingress)
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-app-route
spec:
  parentRefs:
    - name: my-gateway
  hostnames:
    - awesome-http.example.com
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /dev
      backendRefs:
        - name: my-app
          port: 8080
```

## Troubleshooting

### Gateway not getting IP/Port assigned
```bash
kubectl describe gateway my-gateway
kubectl logs -n envoy-gateway-system -l app=envoy-gateway
```

### HTTPRoute not working
```bash
kubectl describe httproute my-app-route
kubectl get httproute -o yaml
```

### Check Gateway Status
```bash
kubectl get gateway my-gateway -o yaml | grep -A 10 status
```

## References

- [Gateway API Official Docs](https://gateway.api.k8s.io/)
- [Envoy Gateway Docs](https://gateway.envoyproxy.io/)
- [Migration Guide from Ingress](https://gateway.api.k8s.io/guides/migrating-from-ingress/)
- [HTTPRoute Spec](https://gateway.api.k8s.io/v1alpha2/api-types/httproute/)

## Next Steps

1. Recreate your cluster using the updated config
2. Choose and install a Gateway controller
3. Convert your existing Ingress rules to Gateway API resources
4. Test connectivity before removing Ingress controller
