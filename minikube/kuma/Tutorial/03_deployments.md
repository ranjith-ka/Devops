# Kuma Deployments

Follow the documents for installting the kuma control plane

```bash
make kuma
```

This will install the kuma in local with only control plan configuration, there is no global zone config settings.

## Zone Ingress

- All requests that are sent from one zone to another will be directed to the proper instance by the Zone Ingress.
- Because ZoneIngress uses Service Name Indication (SNI) to route traffic, mTLS is required to do cross zone communication.

> **Note:** You shouldn't run zoneEgress when running the CP in global



## Zone Egress

ZoneEgress proxy is used when it is required to isolate outgoing traffic (to services in other zones or external services in the local zone). and you want to achieve isolation of outgoing traffic (to services in other zones or external services in the local zone), you can use ZoneEgress proxy.

TODO to test for routing the traffic via Egress

```yaml
echo "apiVersion: kuma.io/v1alpha1
kind: Mesh
metadata:
  name: default
spec:
  routing:
    zoneEgress: true
  mtls: # mTLS is required to use ZoneEgress
    [...]" | kubectl apply -f -
```

### Configure zone proxy authentication

To obtain a configuration from the control plane, a zone proxy (zone ingress / zone egress) must authenticate itself. There are several authentication methods available.

```bash
export ADMIN_TOKEN=$(kubectl get secrets -n kuma-system admin-user-token -ojson | jq -r .data.value | base64 -d)
kumactl config control-planes add --name=kind --headers "authorization=Bearer $ADMIN_TOKEN" --address=http://kuma.example.com --overwrite
```

### Mesh Secrets
