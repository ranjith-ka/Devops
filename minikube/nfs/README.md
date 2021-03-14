# NFS storage for Minikube settings

## Shared Storage

1. NFS used as shared storage for development purpose.
2. HA for NFS not achived in this Research.
3. Deploy NFS chart and use that Storage as Shared storage across pods.
4. Pods can do autoprovisioning with VolumeClaim options.

### Values.yaml for NFS server provisioner

Using NFS charts create the NFS server and NFS volume. Use nfs volumeclaim templates to create voluems under NFS volumes.

Those Volumes created out of NFS server provisiner can be mounted in Pods across.

### Steps

Add the Stable repo in local deployment machine.

`helm repo list`

NAME  URL
stable <https://kubernetes-charts.storage.googleapis.com>

`helm search repo nfs`

| NAME        | CHART VERSION  | APP VERSION | DESCRIPTION |
| ------- | :------- | :-------: |  :-------: |
| stable/nfs-client-provisioner  | `1.2.8` |    `3.1.0`   |    nfs-client is an automatic provisioner that use...|
stable/nfs-server-provisioner |  `1.0.0`   |   `2.3.0`  |     nfs-server-provisioner is an out-of-tree dynami...|

### To install the Chart

`helm install -f minikube/dev/nfs/values.yaml nfs stable/nfs-server-provisioner -n namespace`

Verifiy the Pods. You can see the Statefulset created from this Helm Deployment.

`kubectl get all -n namespace`

`kubectl get sts`

### Troubleshooting

To check the logs in the Pod, while auto provision the NFS volume (it create a PV and PVC as per roles)

`kubectl logs nfs-server-provisioner-0`

`kubectl describe pod nfs-server-provisioner-0`

### Reference

<https://github.com/kubernetes-incubator/external-storage/tree/master/nfs>

<https://github.com/helm/charts/tree/master/stable/nfs-server-provisioner>

