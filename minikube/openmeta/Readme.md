# OpenMeta Data

<https://docs.open-metadata.org>

### Install

Use the instruction to install the openmeta 

<https://docs.open-metadata.org/deployment/kubernetes>

Make sure local kube is running and install required tools like helm, kubectl, make

```bash
make openmeta
```

```bash
kubectl get pod
```

Make changes in the local and validate before moving to CDT.

### Cleanup resources:

```bash
make openmeta-cleanup
```

### Install steps

1. Install the Deps first, so the init scripts will run without error.
2. Local Storage required time to create the PV's (wait untile the storage PV's get created)
  `k describe po helper-pod-create-pvc-12efedaa-d474-4538-b4d6-cc90c367233b -n local-path-storage` 
3. 