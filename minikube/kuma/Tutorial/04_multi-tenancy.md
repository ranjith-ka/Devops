# Configuring your Mesh and multi-tenancy

Multi-tenancy in Kuma, which allows the creation of multiple isolated service meshes within the same Kuma cluster.

This feature simplifies the operation of Kuma in environments where different meshes are required for security, segmentation, or governance reasons.

Mesh is the parent resource in Kuma and includes data plane proxies and policies.

At least one Mesh must exist to use Kuma, and there is no limit to the number of Meshes that can be created. Each data plane proxy can only belong to one Mesh at a time.

## Create a new Mesh

This is created when start from scratch.

```yaml
apiVersion: kuma.io/v1alpha1
kind: Mesh
metadata:
  name: default
```

## Data plane proxy

A data plane proxy (DPP) is the part of Kuma that runs next to each workload that is a member of the mesh. A DPP is composed of the following components:

- Data plane proxies are also called sidecars.

### Inbound

An inbound consists of a set of tags & the port the workload listens on
Most of the time a DPP exposes a single inbound(8080). When a workload exposes multiple ports, multiple inbounds can be defined.

The kuma-dp retrieves Envoy startup configuration from the control plane.
The kuma-dp process starts Envoy with this configuration.
Envoy connects to the control plane using XDS and receives configuration updates when the state of the mesh changes.
The control plane uses policies and Dataplane entities to generate the DPP configuration

Intercepted Traffic:

- Inbound traffic is intercepted by the DPP and forwarded to the Envoy proxy.

<https://kuma.io/docs/2.6.x/production/dp-config/transparent-proxying/#configuration>

Reachable service:

https://kuma.io/docs/2.6.x/production/dp-config/transparent-proxying/#reachable-services


### Secure Deployments <TODO>

Secrets belong to a specific Mesh resource, and cannot be shared across different Meshes. Policies use secrets at runtime.

Kuma leverages Secret resources internally for certain operations, for example when storing auto-generated certificates and keys when Mutual TLS is enabled.

