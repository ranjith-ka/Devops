OS_NAME := $(shell uname -s | tr A-Z a-z)
TIMESTAMP ?= $$(date -u +'%Y%m%d%H%M%S')
PROTOCV := 3.11.4

# strip patch version number
PROTOCVMINOR ?= $(subst $(suffix $(PROTOCV)),,$(PROTOCV))

ifeq ($(findstring mingw64,$(OS_NAME)),mingw64)
PROTOZIP ?= protoc-$(PROTOCV)-win64.zip
PROTOLOC ?= /c/protoc
export PATH := $(PATH):/c/protoc/bin
endif

ifeq ($(findstring linux,$(OS_NAME)),linux)
PROTOZIP ?= protoc-$(PROTOCV)-linux-x86_64.zip
PROTOLOC ?= /usr/local
export PATH := $(PATH):/usr/local/bin
endif

ifeq ($(findstring darwin,$(OS_NAME)),darwin)
PROTOZIP ?= protoc-$(PROTOCV)-osx-x86_64.zip
PROTOLOC ?=  $(HOME)/.protobuf
export PATH := $(PATH):$(HOME)/.protobuf/bin
endif

default: build

install: install-protoc

install-protoc:
	@echo Installing protoc $(PROTOCV) binaries to $(PROTOLOC)
ifeq ($(findstring mingw64,$(OS_NAME)),mingw64)
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@mkdir -p /c/protoc
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc* > /dev/null
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
	@echo
	@echo Please update your PATH to include $(PROTOLOC)/bin
endif

ifeq ($(findstring darwin,$(OS_NAME)),darwin)
	@mkdir -p $(HOME)/.protobuf/bin
	@cd $(HOME)/.protobuf
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc > /dev/null
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
endif

ifeq ($(findstring linux,$(OS_NAME)),linux)
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@sudo unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc > /dev/null
	@sudo unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
endif

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build: generate compile

install-helm:
	@curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

install-kubectl:
	@curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
	@chmod +x kubectl
	@sudo mv kubectl /usr/local/bin/

generate:
	@echo Generate code
	@go generate ./...

compile:
	@echo Compile binary
	@go build .

run: build
	./Devops random

reports:
	@mkdir -p reports

bench: reports
	@echo Writing benchmark to reports/bench$(TIMESTAMP).txt
	@go test -benchmem -run=^$$ <repo> -bench . > reports/bench$(TIMESTAMP).txt

coverage: reports
	@go test -race -covermode=atomic -coverprofile=reports/coverage.out ./...
	@go tool cover -func=reports/coverage.out

dirty:
ifneq ($(DIRTY),)
	@echo modified/untracked files; echo $(DIRTY); exit 1
else
	@echo 'clean'
endif

include make/kind.mk
include make/minikube.mk
include make/mongo.mk

#### Run this command once minikube available
## eval $(minikube -p minikube docker-env)
####

snapshot:
	@echo Build canary docker image
	@docker build --build-arg APP="canary" -t ranjithka/canary:latest .

image:
	@echo Build production docker image
	@docker build --build-arg APP="canary" -t ranjithka/canary:0.0.1 .

push:
	@echo Publish the docker images
	@docker push ranjithka/canary:latest
	@echo Publish Tagged image
	@docker push ranjithka/canary:0.0.1

ingress:
	@helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	@echo Installing Ingress Helm Chart
	@helm install -f minikube/nginx/values.yaml nginx ingress-nginx/ingress-nginx

ingress2:
	@helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	@echo Installing Ingress Helm Chart
	@helm install -f minikube/nginx/values2.yaml nginx2 ingress-nginx/ingress-nginx

install-canary-app:
	@helm install -f minikube/dev/canary.yaml canary-dev charts/dev
	@helm install -f minikube/dev/prd.yaml prd-dev charts/dev

mongo:
	@helm repo add mongodb https://mongodb.github.io/helm-charts
	@helm install community-operator mongodb/community-operator
	@kubectl apply -f minikube/mongo-community-operator/mongo.yaml

monitoring:
	@helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	@helm install -f minikube/prometheus/values.yaml prometheus prometheus-community/prometheus
	@helm repo add grafana https://grafana.github.io/helm-charts
	@helm install -f minikube/grafana/values.yaml grafana grafana/grafana

es:
	@helm install -f minikube/elastic-search/master.yaml elasticsearch-master elastic/elasticsearch
	@helm install -f minikube/elastic-search/data.yaml elasticsearch-data elastic/elasticsearch
	@helm install -f minikube/kibana/values.yaml kibana elastic/kibana

clean-es:
	@kubectl delete pvc elasticsearch-data-elasticsearch-data-0 elasticsearch-master-elasticsearch-master-0

delete-app:
	@helm delete prd-dev canary-dev grafana prometheus

crossplane:
	@kubectl create namespace crossplane-system
	@helm repo add crossplane-stable https://charts.crossplane.io/stable
	@helm install crossplane --namespace crossplane-system crossplane-stable/crossplane

clean-crossplane:
	@helm delete crossplane --namespace crossplane-system
	@kubectl delete namespace crossplane-system
	
vault:
	@helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
	@kubectl create namespace vault-infra && kubectl label namespace vault-infra name=vault-infra
	@helm upgrade --namespace vault-infra --install vault-operator banzaicloud-stable/vault-operator --wait
	@kubectl apply -f minikube/vault/rbac.yaml
	@kubectl apply -f minikube/vault/cr.yaml
	@helm upgrade --namespace vault-infra --install vault-secrets-webhook banzaicloud-stable/vault-secrets-webhook --wait
	
vault-token:
	@kubectl get secrets vault-unseal-keys -o jsonpath={.data.vault-root} | base64 --decode
	@kubectl apply -f minikube/vault/vault-ingress.yaml

delete-vault:
	@kubectl delete -f minikube/vault/rbac.yaml
	@kubectl delete -f minikube/vault/cr.yaml
	@helm delete vault-operator vault-secrets-webhook  -n vault-infra
	@kubectl delete namespace vault-infra

flux:
	@flux install --components-extra=image-reflector-controller,image-automation-controller

openmeta-deps:
	@kubectl create secret generic mysql-secrets --from-literal=openmetadata-mysql-password=openmetadata_password
	@kubectl create secret generic airflow-secrets --from-literal=openmetadata-airflow-password=admin
	@kubectl create secret generic airflow-mysql-secrets --from-literal=airflow-mysql-password=airflow_pass
	@kubectl apply -f minikube/openmeta/source.yaml
	@kubectl apply -f minikube/openmeta/deps_flux_local.yaml

openmeta-cleanup:
	@kubectl delete secrets mysql-secrets airflow-secrets airflow-mysql-secrets
	@kubectl delete -f minikube/openmeta/deps_flux_local.yaml

nfs:
	@helm repo add nfs-ganesha-server-and-external-provisioner https://kubernetes-sigs.github.io/nfs-ganesha-server-and-external-provisioner
	@helm install -f minikube/nfs/values.yaml nfs nfs-ganesha-server-and-external-provisioner/nfs-server-provisioner

kafka:
	@helm repo add strimzi https://strimzi.io/charts/
	@helm install --create-namespace --namespace kafka strimzi strimzi/strimzi-kafka-operator --namespace kafka

kafka-ui:
	@helm repo add kafka-ui https://provectus.github.io/kafka-ui-charts
	@helm install kafka-ui -f minikube/kafka/kafka-ui.yaml kafka-ui/kafka-ui

kuma-global:
	@helm repo add kuma https://kumahq.github.io/charts
	@helm install --create-namespace --namespace kuma-system kuma-global -f minikube/kuma/global.yaml kuma/kuma

kuma-cp:
	@helm repo add kuma https://kumahq.github.io/charts
	@helm install --namespace kuma-system kuma-cp -f minikube/kuma/cp.yaml kuma/kuma

metallb::
	@kubectl get configmap kube-proxy -n kube-system -o yaml | sed -e "s/strictARP: false/strictARP: true/" | kubectl apply -f - -n kube-system
	@kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.3/config/manifests/metallb-native.yaml