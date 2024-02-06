# Service Mesh

Service Mesh is a technology pattern that implements a better way to implement modern networking and connectivity among the different services that make up an application. While it is commonly used in the context of microservices, it can be used to improve connectivity among every architecture and on every platform like VMs and containers

![Arch](https://kuma.io/assets/images/docs/0.4.0/diagram-14.jpg)

* sidecar proxy and sits on the data plane of our requests
* It also creates fragmentation and security issues as more teams try to address the same concerns on different technology stacks.

## Kuma

Kuma helps implement a service mesh approach to distributed deployments as part of the move from monolithic architectures to microservices

* Universal: Kuma supports every environment, including Kubernetes, VMs, and bare metal
* Envoy-based: Kuma uses Envoy as the sidecar proxy, which is an open source proxy designed for cloud native applications
* Multi-zone: Kuma can be deployed in a multi-zone environment to provide high availability and fault tolerance

### Sidecar proxy

* The services delegate all the connectivity and observability concerns to an out-of-process runtime, that will be on the execution path of every request. It will proxy all the outgoing connections and accept all the incoming ones. And of course it will execute traffic policies at runtime, like routing or logging. By using this approach, developers don’t have to worry about connectivity and focus entirely on their services and applications.

`data plane proxy (DP) Vs control plane (CP)`

Data Plane (kuma-dp)

* These are proxies running with your services, handling all mesh traffic. Kuma uses Envoy as its data plane proxy.

Control Plane (kuma-cp)

* This configures the data plane proxies but doesn't interact with mesh traffic directly. Users create policies, and the control plane processes these to configure the data plane proxies

`kuma.io/sidecar-injection: enabled`

* This annotation is used to enable the automatic injection of the Kuma sidecar proxy into the pods of the application. This is done by the Kuma control plane, which is responsible for configuring the sidecar proxies.

### Sizing your control-plane

In short, a control-plane with 4vCPU and 2GB of memory will be able to accommodate more than 1000 data planes.

```yaml
resources:
    requests:
        cpu: 50m
        memory: 64Mi
    limits:
        cpu: 1000m
        memory: 512Mi
```

Check Container patch for increase the resources

<https://kuma.io/docs/2.6.x/production/dp-config/dpp-on-kubernetes/#workload-matching>

### References

<https://kuma.io/docs/2.6.x/introduction/about-service-meshes>