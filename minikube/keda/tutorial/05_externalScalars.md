### External Scalers

https://keda.sh/docs/2.11/scalers/azure-service-bus/

While KEDA ships with a set of built-in scalers, users can also extend KEDA through a GRPC service that implements the same interface as the built-in scalers.

#TODO

To implement external scalar, need to use the GRPC and trigger the driver to accept the scalars.

Will come back and implement this if required.


### Admission Webhooks

- The scaled workload (scaledobject.spec.scaleTargetRef) is already autoscaled by another other sources (other ScaledObject or HPA).
- CPU and/or Memory trigger are used and the scaled workload doesn’t have the requests defined. This rule doesn’t apply to all the workload types, only to Deployment and StatefulSet.
- CPU and/or Memory trigger are the only used triggers and the ScaledObject defines minReplicaCount:0. This rule doesn’t apply to all the workload types, only to Deployment and StatefulSet.