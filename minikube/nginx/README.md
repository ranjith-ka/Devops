# Nginx testing

1. There is nginx running already with SSL offloading in ELB. To test SSL offloading in the nginx ingress.

2. Since most of the customer use Metal LB or on permise LB. This offload SSL in LB or in nginx ingress.

    <https://www.getambassador.io/docs/latest/topics/running/ambassador-with-aws/#aws-load-balancer-notes>

    In Kubernetes, when using the AWS integration and a service of type LoadBalancer, the only types of load balancers that can be created are ELBs and NLBs. When aws-load-balancer-backend-protocol is set to tcp, AWS will create an L4 ELB. When aws-load-balancer-backend-protocol is set to http, AWS will create an L7 ELB.

3. After checking this document, i would not try this NLB in nginx, since we need to route many GRPC service.

4. Need to route the GRPC service in internet with secured way is the end goal.

5. So i conclude to create SSL certificate and offload in nginx proxy/ingress.

## Steps

Add the Stable repo in local deployment machine.

`helm repo list`
`helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx`

NAME   URL
stable <https://kubernetes-charts.storage.googleapis.com>

### Install the nginx chart with version

`helm install -f minikube/nginx/values.yaml nginx ingress-nginx/ingress-nginx`
`helm install -f minikube/nginx/values.yaml nginx ingress-nginx/ingress-nginx --version 3.36.0`

### Kubernetes 1.6+

Follow the below link for kube1.16+

<https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx>
