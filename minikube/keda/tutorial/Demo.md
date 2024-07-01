# Demo on KEDA

Steps:

- GCO devops repo, Helm charts and FluxCD (Short Update)
- KEDA Architecture (Tutorial 01_scaling_intro.md )
- Install KEDA in local kind cluster
- Install a sample app
- Add scaled object to scale the sample app
- Add scaled object to scale the sample app based on Azure Service Bus queue
- In Dev cluster explain the KEDA architecture

```example
 git:(main) ✗ kpo -n keda
NAME                                               READY   STATUS    RESTARTS       AGE
keda-admission-webhooks-99d9c854-t6dsc             1/1     Running   0              2m15s
keda-operator-5f648b87fb-zhwxj                     1/1     Running   1 (116s ago)   2m15s
keda-operator-metrics-apiserver-846455667b-gzdfc   1/1     Running   0              2m15s
```

```example
  git:(main) ✗ k get crds
NAME                                    CREATED AT
clustertriggerauthentications.keda.sh   2023-11-30T04:39:44Z
scaledjobs.keda.sh                      2023-11-30T04:39:44Z
scaledobjects.keda.sh                   2023-11-30T04:39:44Z
triggerauthentications.keda.sh          2023-11-30T04:39:44
```
