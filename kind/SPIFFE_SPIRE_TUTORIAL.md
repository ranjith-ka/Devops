# SPIFFE / SPIRE Tutorial (Kind)

Hands-on guide to the [SPIFFE standard](https://spiffe.io/docs/latest/spiffe-specs/spiffe/) using **SPIRE** (the reference implementation) on your local kind cluster.

**What you will learn**

1. What SPIFFE IDs, SVIDs, and the Workload API are
2. How SPIRE Server + Agent issue identity on Kubernetes
3. Register a workload and fetch an X.509-SVID

Official references:

- [SPIFFE overview (spec)](https://spiffe.io/docs/latest/spiffe-specs/spiffe/)
- [Kubernetes quickstart](https://spiffe.io/docs/latest/try/getting-started-k8s/)
- [spire-tutorials repo](https://github.com/spiffe/spire-tutorials)

---

## 1. SPIFFE in 5 minutes

SPIFFE (**S**ecure **P**roduction **I**dentity **F**ramework **F**or **E**veryone) standardizes **workload identity** so services authenticate each other without IP-based trust or long-lived shared secrets.

It has three core pieces:

| Component | Role |
|-----------|------|
| **SPIFFE ID** | Stable name of a workload, as a URI |
| **SVID** | Cryptographic document that *proves* that ID (passport) |
| **Workload API** | Local API that delivers SVIDs + trust bundles to a process |

### SPIFFE ID

A URI that names an entity:

```text
spiffe://<trust-domain>/<path>
```

Examples:

```text
spiffe://example.org/ns/default/sa/frontend
spiffe://prod.acme.com/payments/api
```

- **Trust domain** — like a DNS zone for identities (e.g. `example.org`). One org / cluster / environment typically has one trust domain.
- **Path** — arbitrary hierarchy you define (often mirrors K8s namespace + service account).

Identity is **who the workload is**, not where it runs (IP/host change; the SPIFFE ID stays stable).

### SVID (SPIFFE Verifiable Identity Document)

An SVID carries the SPIFFE ID and can be verified cryptographically. Document types:

| Type | Common use |
|------|------------|
| **X.509-SVID** | mTLS between services (Envoy, gRPC, etc.) |
| **JWT-SVID** | Auth to APIs / cloud (short-lived bearer) |
| **WIT-SVID** | WebAssembly Identity Token |

Properties required of any SVID:

1. Proven authentic (signed / chain-validated)
2. Proven bound to the presenter

### Workload API

Workloads obtain SVIDs via a **local** SPIFFE Workload API (usually a Unix domain socket). There is **no app-level login token** — the platform (SPIRE Agent) attests the calling process (e.g. K8s pod identity via kubelet) out of band.

The API also returns **trust bundles** (CAs) for your trust domain and federated domains.

```text
┌─────────────┐     Workload API      ┌──────────────┐
│  Workload   │ ◄── (UDS socket) ──── ►│ SPIRE Agent  │
│  (pod/app)  │   X.509 / JWT SVID    │  (per node)  │
└─────────────┘                       └──────┬───────┘
                                             │ attest + enroll
                                      ┌──────▼───────┐
                                      │ SPIRE Server │
                                      │ (CA / policy)│
                                      └──────────────┘
```

**SPIFFE** = the standard. **SPIRE** = open-source server/agent that implements it.

---

## 2. Prerequisites

- kind cluster (see `kind/config.yaml`)
- `kubectl`, `git`
- Optional: `openssl` (to inspect X.509-SVIDs)

Create / recreate the cluster if needed:

```bash
kind create cluster --config kind/config.yaml --name k8s
kubectl cluster-info
kubectl get nodes
```

---

## 3. Clone SPIRE tutorial manifests

```bash
cd /tmp
git clone --depth 1 https://github.com/spiffe/spire-tutorials.git
cd spire-tutorials/k8s/quickstart
```

All `kubectl apply` commands below assume you are in this directory.

---

## 4. Deploy SPIRE Server

### 4.1 Namespace + server RBAC / bundle configmap

```bash
kubectl apply -f spire-namespace.yaml

kubectl apply \
  -f server-account.yaml \
  -f spire-bundle-configmap.yaml \
  -f server-cluster-role.yaml
```

### 4.2 Server config, StatefulSet, Service

```bash
kubectl apply \
  -f server-configmap.yaml \
  -f server-statefulset.yaml \
  -f server-service.yaml
```

Wait until ready:

```bash
kubectl -n spire rollout status statefulset/spire-server
kubectl -n spire get pods,svc
```

Expected: `spire-server-0` **Running**, service `spire-server` exposing port `8081`.

> Kind note: the quickstart StatefulSet uses a PersistentVolumeClaim. Kind’s default storage class usually binds this automatically. If `spire-server-0` stays `Pending` with unbound PVC, check `kubectl -n spire get pvc` and your cluster storage.

---

## 5. Deploy SPIRE Agent

Agents run as a **DaemonSet** (one per node) and talk to kubelet for **workload attestation**.

```bash
kubectl apply \
  -f agent-account.yaml \
  -f agent-cluster-role.yaml

kubectl apply \
  -f agent-configmap.yaml \
  -f agent-daemonset.yaml
```

Verify (your cluster has control-plane + worker → expect **2** agent pods if the DaemonSet schedules on both):

```bash
kubectl -n spire get daemonset,pods
```

---

## 6. Register identities

SPIRE only issues an SVID after a **registration entry** matches selectors for the caller.

### 6.1 Node (Agent) registration

```bash
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
  -selector k8s_psat:cluster:demo-cluster \
  -selector k8s_psat:agent_ns:spire \
  -selector k8s_psat:agent_sa:spire-agent \
  -node
```

This tells the server: agents attested via PSAT with those selectors get SPIFFE ID  
`spiffe://example.org/ns/spire/sa/spire-agent` (and can parent workloads).

### 6.2 Workload registration

Map pods in `default` using SA `default` to a workload ID:

```bash
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://example.org/ns/default/sa/default \
  -parentID spiffe://example.org/ns/spire/sa/spire-agent \
  -selector k8s:ns:default \
  -selector k8s:sa:default
```

List entries:

```bash
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry show
```

---

## 7. Fetch an SVID from a workload

Deploy the sample client (mounts the Agent UDS at `/run/spire/sockets/agent.sock`):

```bash
kubectl apply -f client-deployment.yaml
kubectl wait --for=condition=Ready pod -l app=client --timeout=120s
```

Fetch SVIDs via the Workload API:

```bash
CLIENT_POD=$(kubectl get pods -l app=client -o jsonpath='{.items[0].metadata.name}')

kubectl exec -it "$CLIENT_POD" -- \
  /opt/spire/bin/spire-agent api fetch \
  -socketPath /run/spire/sockets/agent.sock
```

You should see one or more SVIDs, including:

```text
spiffe://example.org/ns/default/sa/default
```

Optional — write certs and inspect the URI SAN:

```bash
kubectl exec "$CLIENT_POD" -- \
  /opt/spire/bin/spire-agent api fetch -write /tmp \
  -socketPath /run/spire/sockets/agent.sock

kubectl exec "$CLIENT_POD" -- \
  openssl x509 -in /tmp/svid.0.pem -noout -text | grep -A2 'Subject Alternative Name'
```

Look for `URI:spiffe://example.org/ns/default/sa/default`.

---

## 8. Mental model: how attestation worked

1. **Agent → Server**: Agent proves it runs as `spire` / SA `spire-agent` on cluster `demo-cluster` (PSAT selectors) → gets node identity.
2. **Workload → Agent**: Client pod mounts `agent.sock`. Agent asks kubelet “who is this process?” → namespace `default`, SA `default`.
3. **Match**: Server registration entry matches those K8s selectors → Agent returns X.509-SVID (and optionally JWT-SVID) + trust bundle.

No secrets baked into the Deployment YAML — identity is **issued at runtime** after attestation.

---

## 9. Map concepts → what you just ran

| Spec concept | In this lab |
|--------------|-------------|
| Trust domain | `example.org` |
| SPIFFE ID | `spiffe://example.org/ns/default/sa/default` |
| X.509-SVID | Cert returned by `spire-agent api fetch` |
| Workload API | UDS `/run/spire/sockets/agent.sock` |
| Trust bundle | CA material delivered with the SVID |
| SPIRE Server | `spire/spire-server-0` |
| SPIRE Agent | DaemonSet `spire-agent` |

---

## 10. Next steps (platform path)

| Goal | Where to look |
|------|----------------|
| mTLS between services | [SPIRE + Envoy + X.509-SVIDs](https://spiffe.io/docs/latest/keyless/envoy/) |
| JWT to cloud / APIs | [JWT-SVIDs](https://spiffe.io/docs/latest/spiffe-specs/jwt-svid/) / [AWS OIDC](https://spiffe.io/docs/latest/keyless/oidc/) |
| Policy on identity | [SPIRE + OPA + Envoy](https://spiffe.io/docs/latest/keyless/opa/) |
| Multi-cluster | [SPIFFE Federation](https://spiffe.io/docs/latest/spiffe-specs/spiffe-federation/) |
| Production Helm | [SPIRE Helm charts](https://spiffe.io/docs/latest/deploying/install-server-agent/) |

For day-2 ops: registration entries should be managed as config (CI/GitOps), not one-off `entry create` in production. Prefer short SVID TTLs; Agent rotates them automatically for healthy workloads.

---

## 11. Tear down

```bash
kubectl delete deployment client --ignore-not-found
kubectl delete namespace spire

kubectl delete clusterrole spire-server-trust-role spire-agent-cluster-role --ignore-not-found
kubectl delete clusterrolebinding spire-server-trust-role-binding spire-agent-cluster-role-binding --ignore-not-found
```

---

## Troubleshooting

| Symptom | Check |
|---------|--------|
| `spire-server-0` Pending | `kubectl -n spire describe pod spire-server-0` — PVC / storage class |
| Agent CrashLoop | Agent logs: `kubectl -n spire logs -l app=spire-agent` — often bundle or server connectivity |
| `fetch` connection refused | Socket mount path; Agent on same node; client volumeMount matches Agent hostPath/socket |
| Empty / wrong SPIFFE ID | `entry show`; selectors `k8s:ns` / `k8s:sa` must match the pod |
| PSAT / node entry fails | Cluster name in selectors must match Agent config (`demo-cluster` in the quickstart) |
