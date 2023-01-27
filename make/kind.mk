.PHONY: charts

cluster: kind-cluster

## Start Colima
colima:
	@echo #### starting Colima with 4CPU and 8GB Memory ####
	@colima start --cpu 4 --memory 8

kind: colima kind-cluster

kind-cluster:
	@echo Creating Kind environment
	@kind create cluster --config kind/config.yaml --name k8s

load-image:
	@kind load docker-image  ranjithka/canary:0.0.1  --name k8s
	@kind load docker-image  ranjithka/canary:latest --name k8s

delete-kind:
	@kind delete cluster --name k8s
	@colima stop