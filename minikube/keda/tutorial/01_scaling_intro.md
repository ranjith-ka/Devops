# KEDA concepts

https://keda.sh/docs/2.11/concepts/scaling-deployments/#scaling-of-deployments-and-statefulsets

KEDA is a Kubernetes-based Event Driven Autoscaler. With KEDA, you can drive the scaling of any container in Kubernetes based on the number of events needing to be processed.

- With KEDA you can explicitly map the apps you want to use event-driven scale, with other apps continuing to function.

Agent - keda-operator runs with keda install
Metrics -  The metric serving is the primary role of the keda-operator-metrics-apiserver container that runs when you install KEDA.
Admission Webhooks - it will prevent multiple ScaledObjects to target the same scale target.


![Architecture](https://keda.sh/img/keda-arch.png)

#TODO Try to explain this arch again

### Event sources and scalers

- CPU/Memory
- MSSQL
- PostgreSQL etc...

### Custom Resources (CRD)

- scaledobjects.keda.sh
- scaledjobs.keda.sh (Not used)
- triggerauthentications.keda.sh
- clustertriggerauthentications.keda.sh (not used)

ScaledObjects represent the desired mapping between an event source (e.g. Rabbit MQ) and the Kubernetes Deployment, StatefulSet or any Custom Resource that defines /scale subresource.

ScaledJobs represent the mapping between event source and Kubernetes Job.

ScaledObject/ScaledJob may also reference a TriggerAuthentication or ClusterTriggerAuthentication which contains the authentication configuration or secrets to monitor the event source.

Authentication scopes: Namespace vs. Cluster  (ClusterTriggerAuthentication)

- Each TriggerAuthentication is defined in one namespace and can only be used by a ScaledObject in that same namespace.

- For cases where you want to share a single set of credentials between scalers in many namespaces, you can instead create a ClusterTriggerAuthentication.

- As a global object, this can be used from any namespace.

- To set a trigger to use a ClusterTriggerAuthentication, add a kind field to the authentication reference:
