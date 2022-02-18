.PHONY: charts

cluster: kind-cluster

kind: kind-cluster

kind-cluster:
	@echo Creating Kind environment
	@kind create cluster --config kind/config.yaml --name k8s-1.21.1

load-image:
	@kind load docker-image  ranjithka/canary:0.0.1  --name k8s-1.21.1
	@kind load docker-image  ranjithka/prd:0.0.1  --name k8s-1.21.1

delete-kind:
	@kind delete cluster --name k8s-1.21.1
