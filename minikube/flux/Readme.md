# Flux 

Kudos to flux development team, making open source and lot of cool stuffs.

Please read the tutorials and practise each stage and come here to deploy the one step deployments. Flux team developed many days to reach here, add effort to understand before one step to deploy.

1. Run below to run the application in local with canary deployments

```bash
$ make kind ## Note: This will start Colima and this required some prerequisite like kind, kubectl, etc...
$ flux install ## Simple way to start the cluster in local, this repo is not for Platform Engineers to adopt, this is for devops and dev guys.
$ kubectl apply -f minikube/flux/staging  ## Run this to install the app in cluster with Gitops management
```

2. Once the application is running, Gitops to manage the application.
    - Edit the values.yaml and push the changes to remote, forgot the changes for 10mins
    - check in the cluster this changes are applied ? Should be there, also do a sync to see changes immediate.

```bash
$ flux reconcile source git devops
```

Note: If you change the Tag: latest in helm chart, `helm ls` will give the version from chart.yaml only, `kubectl describe po name` will give which sha and image tag is applied.

```bash
➜  ~/code/github/ranjith/Devops git:(main) ✗ helm ls
NAME      	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART                   	APP VERSION
canary-dev	default  	1       	2023-01-29 18:47:44.051134198 +0000 UTC	deployed	dev-3.2.1               	0.0.1
nginx     	default  	1       	2023-01-29 18:47:45.267912657 +0000 UTC	deployed	ingress-nginx-4.0.13    	1.1.0
prd-dev   	default  	2       	2023-01-29 18:51:10.533669879 +0000 UTC	deployed	dev-3.2.1+91a457debe80.1	0.0.1
```