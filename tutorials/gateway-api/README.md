# Gateway API Tutorial - Complete Hands-On Guide

This tutorial will walk you through installing, configuring, and using Kubernetes Gateway API on a local kind cluster.

**Table of Contents:**
1. [Prerequisites](#prerequisites)
2. [Part 1: Setup & Installation](#part-1-setup--installation)
3. [Part 2: Basic HTTP Routing](#part-2-basic-http-routing)
4. [Part 3: Advanced Routing](#part-3-advanced-routing)
5. [Part 4: TLS/HTTPS](#part-4-tlshttps)
6. [Part 5: Cross-Namespace Routing](#part-5-cross-namespace-routing)
7. [Part 6: Traffic Management](#part-6-traffic-management)
8. [Part 7: Troubleshooting](#part-7-troubleshooting)
9. [Migration from Ingress](#migration-from-ingress)

---

## Prerequisites

- Docker installed and running
- kind CLI installed: `brew install kind`
- kubectl installed: `brew install kubernetes-cli`
- helm installed: `brew install helm`
- Basic understanding of Kubernetes concepts

Verify your setup:
```bash
kind version
kubectl version --client
helm version
```

---

## Part 1: Setup & Installation

### Step 1.1: Create Kind Cluster with Gateway API Support

Your updated `kind/config.yaml` already has Gateway API feature gates enabled. Let's create the cluster:

```bash
# Navigate to the Devops directory
cd /Users/ranjith.a/code/github/ranjith/Devops

# Delete existing cluster if needed
kind delete cluster --name k8s

# Create new cluster with Gateway API support
kind create cluster --config kind/config.yaml --name k8s

# Verify cluster is running
kubectl cluster-info
kubectl get nodes
```

Expected output:
```
kubernetes.io/hostname=k8s-control-plane
kubernetes.io/hostname=k8s-worker
```

### Step 1.2: Install Gateway API Custom Resource Definitions (CRDs)

The CRDs define the new resource types we'll use (Gateway, HTTPRoute, etc.):

```bash
# Install stable Gateway API CRDs (v1.2.1)
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# Verify CRDs are installed
kubectl get crds | grep gateway.networking.k8s.io
```

Expected CRDs:
```
backendtlspolicies.gateway.networking.k8s.io
gatewayclasses.gateway.networking.k8s.io
gateways.gateway.networking.k8s.io
grpcroutes.gateway.networking.k8s.io
httproutes.gateway.networking.k8s.io
referencepolicies.gateway.networking.k8s.io
tcproutes.gateway.networking.k8s.io
tlsroutes.gateway.networking.k8s.io
udproutes.gateway.networking.k8s.io
```

### Step 1.3: Install Envoy Gateway Controller

The Gateway CRD alone doesn't do anything—we need a controller that implements it. We'll use Envoy Gateway:

```bash
# Add the Envoy Gateway Helm repository
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo update

# Install Envoy Gateway in the envoy-gateway-system namespace
helm install envoy-gateway envoy-gateway/gateway-helm \
  --namespace envoy-gateway-system \
  --create-namespace \
  --values - << 'EOF'
config:
  envoyGateway:
    logging:
      level:
        default: debug
EOF

# Wait for the controller to be ready (30-60 seconds)
kubectl wait --timeout=5m -n envoy-gateway-system \
  --for=condition=Progressing=True \
  deployment/envoy-gateway

# Verify installation
kubectl get pods -n envoy-gateway-system
kubectl get gatewayclasses
```

Expected `gatewayclasses` output:
```
NAME              CONTROLLER                                      ACCEPTED   AGE
envoy             gateway.envoyproxy.io/gatewayclass-controller   True       1m
```

### Step 1.4: Verify Gateway API is Ready

```bash
# Check if the Envoy Gateway controller is running
kubectl logs -n envoy-gateway-system -l app=envoy-gateway --tail=20

# Check available APIs
kubectl api-resources | grep gateway
```

✅ **Setup Complete!** You now have Gateway API and Envoy Gateway controller running.

---

## Part 2: Basic HTTP Routing

### Step 2.1: Deploy a Sample Application

First, let's create a simple web application to test routing:

```bash
# Create a namespace for the app
kubectl create namespace demo-apps

# Deploy the app
kubectl apply -f - << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-app
  namespace: demo-apps
  labels:
    app: http-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: http-app
  template:
    metadata:
      labels:
        app: http-app
    spec:
      containers:
      - name: app
        image: hashicorp/http-echo:latest
        args:
          - "-listen=:8080"
          - "-text=Hello from HTTP App! 🚀"
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: http-app
  namespace: demo-apps
spec:
  selector:
    app: http-app
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
EOF

# Verify deployment
kubectl get pods -n demo-apps
kubectl get svc -n demo-apps
```

### Step 2.2: Create Your First Gateway

A Gateway represents a load balancer and defines the listeners (ports/protocols):

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: demo-gateway
  namespace: demo-apps
spec:
  gatewayClassName: envoy
  listeners:
  - name: http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same
EOF

# Check Gateway status
kubectl get gateway -n demo-apps
kubectl describe gateway demo-gateway -n demo-apps
```

Look for the "Address" field in the Gateway status. It might take a moment to get an address.

### Step 2.3: Create an HTTPRoute

HTTPRoute defines the actual routing rules:

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: http-app-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "demo.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
      weight: 100
EOF

# Check HTTPRoute status
kubectl get httproute -n demo-apps
kubectl describe httproute http-app-route -n demo-apps
```

### Step 2.4: Test the Route

Find the gateway's IP and port:

```bash
# Get the gateway details
kubectl get gateway demo-gateway -n demo-apps -o wide

# Forward the port to your local machine (for testing)
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

# Test with curl
curl -H "Host: demo.localhost" http://localhost:8080/

# You should see:
# Hello from HTTP App! 🚀
```

Stop the port-forward when done:
```bash
pkill -f "port-forward"
```

✅ **Basic HTTP Routing Complete!** You've created a working Gateway API setup.

---

## Part 3: Advanced Routing

### Step 3.1: Multiple Paths Routing

Route different paths to different services. First, deploy another app:

```bash
kubectl apply -f - << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-app
  namespace: demo-apps
spec:
  replicas: 1
  selector:
    matchLabels:
      app: api-app
  template:
    metadata:
      labels:
        app: api-app
    spec:
      containers:
      - name: app
        image: hashicorp/http-echo:latest
        args:
          - "-listen=:8080"
          - "-text=API Response - v1"
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: api-app
  namespace: demo-apps
spec:
  selector:
    app: api-app
  ports:
  - port: 8080
    targetPort: 8080
EOF

# Update HTTPRoute to route different paths
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: http-app-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "demo.localhost"
  rules:
  # Rule 1: /api/* goes to api-app
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: api-app
      port: 8080
  # Rule 2: /* goes to http-app (default)
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
EOF

# Test
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

curl -H "Host: demo.localhost" http://localhost:8080/
# Output: Hello from HTTP App! 🚀

curl -H "Host: demo.localhost" http://localhost:8080/api
# Output: API Response - v1

pkill -f "port-forward"
```

### Step 3.2: Header-Based Routing

Route based on request headers:

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: header-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "demo.localhost"
  rules:
  # Route to api-app if X-Version header = v2
  - matches:
    - headers:
      - name: X-Version
        value: v2
      path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: api-app
      port: 8080
  # Default route
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
EOF

# Test header routing
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

curl -H "Host: demo.localhost" http://localhost:8080/api
# Output: API Response - v1

curl -H "Host: demo.localhost" -H "X-Version: v2" http://localhost:8080/api
# Same output (both route to api-app)

pkill -f "port-forward"
```

### Step 3.3: Multiple Hostnames

Route different hosts to different services:

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: multi-host-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "app.localhost"
  - "api.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: api-host-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "api.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: api-app
      port: 8080
EOF

# Test
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

curl -H "Host: app.localhost" http://localhost:8080/
# Output: Hello from HTTP App! 🚀

curl -H "Host: api.localhost" http://localhost:8080/
# Output: API Response - v1

pkill -f "port-forward"
```

✅ **Advanced Routing Complete!**

---

## Part 4: TLS/HTTPS

### Step 4.1: Create Self-Signed Certificate

For local testing with TLS:

```bash
# Create a cert directory
mkdir -p /tmp/certs && cd /tmp/certs

# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/CN=demo.localhost"

# Create Kubernetes secret
kubectl create secret tls demo-tls \
  --cert=/tmp/certs/cert.pem \
  --key=/tmp/certs/key.pem \
  -n demo-apps

# Verify secret
kubectl get secret -n demo-apps
```

### Step 4.2: Update Gateway to Support HTTPS

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: demo-gateway
  namespace: demo-apps
spec:
  gatewayClassName: envoy
  listeners:
  # HTTP listener
  - name: http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same
  # HTTPS listener
  - name: https
    port: 443
    protocol: HTTPS
    tls:
      mode: Terminate
      certificateRefs:
      - name: demo-tls
        kind: Secret
    allowedRoutes:
      namespaces:
        from: Same
EOF

# Verify gateway
kubectl describe gateway demo-gateway -n demo-apps
```

### Step 4.3: Update HTTPRoute for HTTPS

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: https-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
    sectionName: https
  hostnames:
  - "demo.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
EOF

# Test HTTPS (insecure for self-signed cert)
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8443:443 &

curl -k -H "Host: demo.localhost" https://localhost:8443/
# Output: Hello from HTTP App! 🚀

pkill -f "port-forward"
```

✅ **TLS/HTTPS Complete!**

---

## Part 5: Cross-Namespace Routing

### Step 5.1: Create Another Namespace with App

```bash
# Create namespace
kubectl create namespace api-backend

# Deploy app
kubectl apply -f - << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-app
  namespace: api-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend-app
  template:
    metadata:
      labels:
        app: backend-app
    spec:
      containers:
      - name: app
        image: hashicorp/http-echo:latest
        args:
          - "-listen=:8080"
          - "-text=Backend Service Response"
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: backend-app
  namespace: api-backend
spec:
  selector:
    app: backend-app
  ports:
  - port: 8080
    targetPort: 8080
EOF
```

### Step 5.2: Create ReferencePolicy for Cross-Namespace Access

By default, a route cannot reference a service in a different namespace. We need a ReferencePolicy:

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferencePolicy
metadata:
  name: allow-cross-namespace
  namespace: api-backend
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: demo-apps
  to:
  - group: ""
    kind: Service
    name: backend-app
EOF
```

### Step 5.3: Create HTTPRoute Referencing Cross-Namespace Service

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: backend-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "backend.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: backend-app
      namespace: api-backend
      port: 8080
EOF

# Test
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

curl -H "Host: backend.localhost" http://localhost:8080/
# Output: Backend Service Response

pkill -f "port-forward"
```

✅ **Cross-Namespace Routing Complete!**

---

## Part 6: Traffic Management

### Step 6.1: Weighted Load Balancing

Route traffic to multiple services with different weights:

```bash
# Deploy a second version of http-app
kubectl apply -f - << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-app-v2
  namespace: demo-apps
spec:
  replicas: 1
  selector:
    matchLabels:
      app: http-app-v2
  template:
    metadata:
      labels:
        app: http-app-v2
    spec:
      containers:
      - name: app
        image: hashicorp/http-echo:latest
        args:
          - "-listen=:8080"
          - "-text=Hello from HTTP App V2! 🚀"
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: http-app-v2
  namespace: demo-apps
spec:
  selector:
    app: http-app-v2
  ports:
  - port: 8080
    targetPort: 8080
EOF

# Create HTTPRoute with weighted traffic
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: weighted-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "weighted.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    # 70% of traffic goes to v1
    - name: http-app
      port: 8080
      weight: 70
    # 30% of traffic goes to v2
    - name: http-app-v2
      port: 8080
      weight: 30
EOF

# Test (make multiple requests to see both versions)
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80 &

for i in {1..10}; do
  curl -s -H "Host: weighted.localhost" http://localhost:8080/
  echo " (Request $i)"
done
# You should see a mix of v1 and v2 responses (70/30 split)

pkill -f "port-forward"
```

### Step 6.2: Request Mirroring

Mirror traffic to another service for testing:

```bash
kubectl apply -f - << 'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: mirror-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - "mirror.localhost"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
    # Mirror requests to v2 for testing (non-blocking)
    filters:
    - type: RequestMirror
      requestMirror:
        backendRef:
          name: http-app-v2
          port: 8080
EOF

# Note: Request mirroring support depends on the controller implementation
# Envoy Gateway may require additional configuration
```

✅ **Traffic Management Complete!**

---

## Part 7: Troubleshooting

### Common Issues and Solutions

#### Issue 1: Gateway Not Getting an Address

```bash
# Check gateway status
kubectl describe gateway demo-gateway -n demo-apps

# Check controller logs
kubectl logs -n envoy-gateway-system -l app=envoy-gateway --tail=50

# Check if controller is running
kubectl get pods -n envoy-gateway-system
```

**Solution:** Wait for the Envoy Gateway controller to be fully ready, or check for resource constraints.

#### Issue 2: HTTPRoute Not Routing Traffic

```bash
# Check route status
kubectl describe httproute http-app-route -n demo-apps

# Check for conditions and errors
kubectl get httproute -n demo-apps -o yaml | grep -A 10 status

# Check backend service exists
kubectl get svc -n demo-apps
kubectl get pods -n demo-apps
```

**Solution:** Ensure services and pods are running, and backend references are correct.

#### Issue 3: Debugging Request Flow

```bash
# Enable debug logging in Envoy Gateway
kubectl edit deployment envoy-gateway -n envoy-gateway-system
# Add --log-level=debug flag

# Check Envoy proxy logs
kubectl logs -n envoy-gateway-system -l app.kubernetes.io/name=envoy --tail=50

# Verify DNS resolution
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup http-app.demo-apps.svc.cluster.local
```

#### Issue 4: Port-Forward Connection Refused

```bash
# Ensure port-forward command is correct
kubectl port-forward -n demo-apps svc/envoy-gateway-demo-gateway 8080:80

# In another terminal, test
curl -H "Host: demo.localhost" http://localhost:8080/

# Check if process is running
ps aux | grep port-forward
```

### Useful Debugging Commands

```bash
# Get all Gateway resources
kubectl get gateways,httproutes,grpcroutes -A

# Watch resource changes
kubectl get httproutes -n demo-apps -w

# Export resource YAML for inspection
kubectl get httproute http-app-route -n demo-apps -o yaml

# Check resource events
kubectl describe httproute http-app-route -n demo-apps

# Check controller configuration
kubectl get gatewayclass envoy -o yaml

# View all CRDs related to Gateway API
kubectl api-resources | grep gateway
```

---

## Migration from Ingress

### Before (Ingress)

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: demo-ingress
  namespace: demo-apps
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/rate-limit: "10"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - demo.localhost
    secretName: demo-tls
  rules:
  - host: demo.localhost
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: http-app
            port:
              number: 8080
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-app
            port:
              number: 8080
```

### After (Gateway API)

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: demo-gateway
  namespace: demo-apps
spec:
  gatewayClassName: envoy
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
      - name: demo-tls
        kind: Secret
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: demo-route
  namespace: demo-apps
spec:
  parentRefs:
  - name: demo-gateway
    kind: Gateway
  hostnames:
  - demo.localhost
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: api-app
      port: 8080
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: http-app
      port: 8080
```

### Key Differences

| Feature | Ingress | Gateway API |
|---------|---------|------------|
| **Port Configuration** | Limited; implicit | Explicit listeners |
| **TLS** | Via annotations | Native in spec |
| **Routing Flexibility** | Path/Host only | Headers, methods, params, etc |
| **Rules Order** | First match | All matches evaluated |
| **Cross-namespace** | Via class | Via ReferencePolicy |
| **Rate Limiting** | Via annotations | Via policy resources |
| **Resource Separation** | Single resource | Gateway + Routes |

---

## Cleanup

When you're done, clean up resources:

```bash
# Delete namespaces (cascades delete resources)
kubectl delete namespace demo-apps api-backend

# Uninstall Envoy Gateway
helm uninstall envoy-gateway -n envoy-gateway-system
kubectl delete namespace envoy-gateway-system

# Delete CRDs (optional, if you want to fully remove Gateway API)
kubectl delete crd -l gateway.networking.k8s.io/

# Delete kind cluster (optional)
kind delete cluster --name k8s
```

---

## Summary

You've learned:
✅ Setting up Gateway API and Envoy Gateway  
✅ Creating Gateways and HTTPRoutes  
✅ Basic and advanced routing patterns  
✅ TLS/HTTPS termination  
✅ Cross-namespace routing with ReferencePolicy  
✅ Traffic management (weighted load balancing)  
✅ Troubleshooting common issues  
✅ Migrating from Ingress to Gateway API  

---

## Additional Resources

- [Official Gateway API Docs](https://gateway.api.k8s.io/)
- [Envoy Gateway Documentation](https://gateway.envoyproxy.io/)
- [Gateway API Concepts](https://gateway.api.k8s.io/concepts/api-overview/)
- [HTTPRoute Specification](https://gateway.api.k8s.io/v1alpha2/api-types/httproute/)
- [Envoy Proxy Documentation](https://www.envoyproxy.io/docs/)

---

## Next Steps

1. Explore other route types: `GRPCRoute`, `TCPRoute`, `TLSRoute`
2. Implement custom Gateway controllers (Cilium, Kong, Nginx)
3. Set up observability with Prometheus and Grafana
4. Implement security policies with Kubernetes Network Policies
5. Deploy in production clusters
