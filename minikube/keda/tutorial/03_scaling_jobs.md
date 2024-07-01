### Scaling Jobs

The primary reason to consider this option is to handle processing long-running executions. Rather than processing multiple events within a deployment, for each detected event a single Kubernetes Job is scheduled. That job will initialize, pull a single event from the message source, and process to completion and terminate.


```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledJob
metadata:
  name: {scaled-job-name}
  annotations:
    autoscaling.keda.sh/paused: true          # Optional. Use to pause autoscaling of Jobs
spec:
  jobTargetRef:
    parallelism: 1                            # [max number of desired pods](https://kubernetes.io/docs/concepts/workloads/controllers/job/#controlling-parallelism)
    completions: 1                            # [desired number of successfully finished pods](https://kubernetes.io/docs/concepts/workloads/controllers/job/#controlling-parallelism)
    activeDeadlineSeconds: 600                #  Specifies the duration in seconds relative to the startTime that the job may be active before the system tries to terminate it; value must be positive integer
    backoffLimit: 6                           # Specifies the number of retries before marking this job failed. Defaults to 6
    template:
      # describes the [job template](https://kubernetes.io/docs/concepts/workloads/controllers/job)
  pollingInterval: 30                         # Optional. Default: 30 seconds
  successfulJobsHistoryLimit: 5               # Optional. Default: 100. How many completed jobs should be kept.
  failedJobsHistoryLimit: 5                   # Optional. Default: 100. How many failed jobs should be kept.
  envSourceContainerName: {container-name}    # Optional. Default: .spec.JobTargetRef.template.spec.containers[0]
  minReplicaCount: 10                         # Optional. Default: 0
  maxReplicaCount: 100                        # Optional. Default: 100
  rolloutStrategy: gradual                    # Deprecated: Use rollout.strategy instead (see below).
  rollout:
    strategy: gradual                         # Optional. Default: default. Which Rollout Strategy KEDA will use.
    propagationPolicy: foreground             # Optional. Default: background. Kubernetes propagation policy for cleaning up existing jobs during rollout.
  scalingStrategy:
    strategy: "custom"                        # Optional. Default: default. Which Scaling Strategy to use. 
    customScalingQueueLengthDeduction: 1      # Optional. A parameter to optimize custom ScalingStrategy.
    customScalingRunningJobPercentage: "0.5"  # Optional. A parameter to optimize custom ScalingStrategy.
    pendingPodConditions:                     # Optional. A parameter to calculate pending job count per the specified pod conditions
      - "Ready"
      - "PodScheduled"
      - "AnyOtherCustomPodCondition"
    multipleScalersCalculation : "max" # Optional. Default: max. Specifies how to calculate the target metrics when multiple scalers are defined.
  triggers:
  # {list of triggers to create jobs}
```

 envSourceContainerName:
  - This optional property specifies the name of container in the Job, from which KEDA should try to get environment properties holding secrets etc. If it is not defined it, KEDA will try to get environment properties from the first Container, ie. from .spec.JobTargetRef.template.spec.containers[0].