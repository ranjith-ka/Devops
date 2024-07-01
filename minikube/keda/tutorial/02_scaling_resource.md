## Scaling Deployments, StatefulSets & Custom Resources

- Deployments and StatefulSets are the most common way to scale workloads with KEDA.
- It allows you to define the Kubernetes Deployment or StatefulSet that you want KEDA to scale based on a scale trigger
- KEDA will monitor that service(Deployment/Statefullset) and based on the events that occur it will automatically scale your resource out/in accordingly.

### ScaledObject spec

https://keda.sh/docs/2.11/concepts/scaling-deployments/#scaledobject-spec

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: {scaled-object-name}
  annotations:
    scaledobject.keda.sh/transfer-hpa-ownership: "true"      # Optional. Use to transfer an existing HPA ownership to this ScaledObject
    autoscaling.keda.sh/paused-replicas: "0"                # Optional. Use to pause autoscaling of objects
spec:
  scaleTargetRef:
    apiVersion:    {api-version-of-target-resource}         # Optional. Default: apps/v1
    kind:          {kind-of-target-resource}                # Optional. Default: Deployment
    name:          {name-of-target-resource}                # Mandatory. Must be in the same namespace as the ScaledObject
    envSourceContainerName: {container-name}                # Optional. Default: .spec.template.spec.containers[0]
  pollingInterval:  30                                      # Optional. Default: 30 seconds
  cooldownPeriod:   300                                     # Optional. Default: 300 seconds
  idleReplicaCount: 0                                       # Optional. Default: ignored, must be less than minReplicaCount 
  minReplicaCount:  1                                       # Optional. Default: 0
  maxReplicaCount:  100                                     # Optional. Default: 100
  fallback:                                                 # Optional. Section to specify fallback options
    failureThreshold: 3                                     # Mandatory if fallback section is included
    replicas: 6                                             # Mandatory if fallback section is included
  advanced:                                                 # Optional. Section to specify advanced options
    restoreToOriginalReplicaCount: true/false               # Optional. Default: false
    horizontalPodAutoscalerConfig:                          # Optional. Section to specify HPA related options
      name: {name-of-hpa-resource}                          # Optional. Default: keda-hpa-{scaled-object-name}
      behavior:                                             # Optional. Use to modify HPA's scaling behavior
        scaleDown:
          stabilizationWindowSeconds: 300
          policies:
          - type: Percent
            value: 100
            periodSeconds: 15
  triggers:
  # {list of triggers to activate scaling of the target resource}
```

```yaml
spec:
  scaleTargetRef:
    apiVersion:    {api-version-of-target-resource}         # Optional. Default: apps/v1
    kind:          {kind-of-target-resource}                # Optional. Default: Deployment
    name:          {name-of-target-resource}                # Mandatory. Must be in the same namespace as the ScaledObject
    envSourceContainerName: {container-name}                # Optional. Default: .spec.template.spec.containers[0]
  ```


  https://keda.sh/docs/2.11/concepts/scaling-deployments/#fallback

  ### fallback

  - The fallback section is optional. It defines a number of replicas to fall back to if a scaler is in an error state.


  ### Triggers

  https://keda.sh/docs/2.11/scalers/

   - type
   - metadata
   - name
   - useCachedMetrics
   - authenticationRef
   - metricType: AverageValues, Value, Utlization


### Caching Metrics

Polling Interval -> HPA -> Metrics -> KEDA Metrics server

- Enabling this feature can significantly reduce the load on the scaler service.

### Pause AutoScaling

```yaml
metadata:
  annotations:
    autoscaling.keda.sh/paused-replicas: "0"
```

### Activating and Scaling thresholds 

- Activation Phase: Defines when the scaler is active or not and scales from/to 0 based on it.
- Scaling Phase: 1 to n vice versa

### Existing HPA

```yaml
metadata:
  annotations:
    scaledobject.keda.sh/transfer-hpa-ownership: "true"
spec:
   advanced:
      horizontalPodAutoscalerConfig:
      name: {name-of-hpa-resource}
```