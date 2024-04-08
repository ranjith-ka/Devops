# Kuma in Prodction

<https://kuma.io/docs/2.6.x/production/>

    -  Deployment topologies
    -  single-zone (one control plane (that can be scaled horizontally) and many data planes connecting directly to it)
    -  multi-zone (advanced deployment model to support multiple Kubernetes or VM-based zones, or hybrid Service Meshes running on both Kubernetes and VMs combined)

## Control plane and data plane architecture

Use kumactl to configure a multi-zone or single-zone control plane, depending on your organization’s needs. You can deploy either a Kubernetes or Universal data plane.

## Single-zone deployment

    - Zone control plane (kuma-cp)
    - Data plane proxies (kuma-dp)
    - Service Connectivity (service - service (via DP))

### zone control plane

    - Accept connections from data plane proxies.
    - Handle policies -> apply to dataplane (in inventory)
    - Manage the config of XDS to DP

### Data plane proxies

    - Connect to other data plane proxies and control plane.

## Multi-zone deployment

A zone can be a Kubernetes cluster, a VPC, or any other deployment you need to include in the same distributed mesh environment. The only condition is that all the data planes running within the zone must be able to connect to the other data planes in this same zone.

    Zone-a & Zone-b (with or without Egress)

    service -> service (via Zone Ingress), This ZoneIngress resource is then also synchronized to the global control plane.

The global control-plane will propagate the zone ingress resources and all policies to all other zones over Kuma Discovery Service (KDS), which is a protocol based on xDS.

![alt text](https://kuma.io/assets/images/diagrams/gslides/kuma_multizone_without_egress.svg)

### Components of a multi-zone deployment

    - Global control plane (kuma-cp)
    - Zone control plane (kuma-cp)
    - Data plane proxies (kuma-dp)
    - Zone Ingress (kuma-ingress)
    - Zone Egress (kuma-egress)

### Failure modes

    - Global control plane failure
        - No policies updates
        - new service will not be discoverable in other zones
        - service removed still appears
        - Zone deletion not possible


    - Zone control plane failure
        - New data plane proxies won't be able to join the mesh. This includes new instances (Pod/VM) that are newly created by automatic deployment mechanisms (e.g., rolling update process), meaning a control plane failure will prevent new instances from joining the mesh.
        - mTLS certiface refresh will fail
        - Data plane proxy configuration will not be updated

### Communication between Global and Zone control plane failing

### Communication between 2 zones failing

<! Too much info before i understand this prodcut>

## kumactl

You can configure kumactl to point to any zone kuma-cp instance by running:
    ```
    $ kumactl config control-planes add --name=XYZ --address=http://{address-to-kuma}:5681
    $ kumactl get meshes
    ```

Kuma - being an application that wants to improve the underlying connectivity between your services by making the underlying network more reliable - also comes with some networking requirements itself.

## Control plane ports

    5682: HTTPS version of the services available under 5681
    5683: gRPC Intercommunication CP server used internally by Kuma to communicate between CP instances.
    5685: the Kuma Discovery Service port, leveraged in multi-zone deployments.

