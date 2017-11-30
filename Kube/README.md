### Adding Readme for my reference. 

To access minikube url 

`minikube service jenkinsv2 --url`

```http://192.168.99.100:30000```

# Deployments 

`kubectl get deployments`
`kubectl get rs`


# Replication controller for Stateless applications
`kubectl scale --replicas=4 -f replication-controller/replicationclr.yml`

# Expose the NodePort  
`kubectl expose deployment nginx-deployment --type=NodePort`

`kubectl rollout status deployment nginx-deployment`

`minikube service nginx-deployment --url`

# To record the deployments 
`kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1 --record`

# Roll back previos version 
`kubectl rollout undo deployment.v1.apps/nginx-deployment`

# Run busybox 
`kubectl run -i --tty busybox --image=busybox --restart=Never -- /bin/sh` 


# Create Config map
`kubectl create configmap nginx-config --from-file=configmap/reverseproxy.conf`

`kubectl list configmap`
`kubectl list secrets`

# Podpreset
