.PHONY: minikube

KUBE_VERSION := v1.23.1
CPU := 4  ## Half of the CPU is good, more than half is bad for my MAC, may be good for a linux servers.
DISK := 40g

### Still testing this, ingress i need to map the port 32080 with localhost 80, So the ingress not working at the moment.
minikube:
	@echo Starting minikube with kube version 
	@minikube start --kubernetes-version=$(KUBE_VERSION) --cpus $(CPU) --disk-size=$(DISK)

#### Run this command once minikube available
## eval $(minikube -p minikube docker-env)
####

mapp: ingress snapshot install-app monitoring