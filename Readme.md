#Helm is an application package manager for Kubernetes, which coordinates the download, installation, and deployment of apps.Helm charts are the way we can define an application as a collection of related Kubernetes resources. 
#create a new helm charts
# $ helm create vizualplatform
# vizualplatform is the chart name
# $ ls
# you will see the below following list
# chart.yaml values.yaml templates charts
# Chart.yaml and values.yaml—define what the chart is and what values will be in it at deployment.
# The chart.yaml application charts are a collection of templates that can be packaged into versioned archives to be deployed.
# The most important part of the chart is the template directory. It holds all the configurations for your application that will be deployed into the cluster
# The values.yaml file you can pass the values following values
# repository where you are pulling your image and the pullPolicy.
# image:
    repository: nginx
    pullPolicy: Always
# Always means it will pull the image on every deployment
# imagePullSecrets: []
  nameOverride: "vizualplatform-app"
  fullnameOverride:"vizualplatform-chart"
# if you need to rename a chart after you create it u can use nameoverride and fullnameoverride.helpers.tpl file is used.
# service accounts name "vizualplatform"
# serviceAccount:
    create: true
    annotations: {}
    name: "vizualplatform"
# NodePort, which exposes the service on each Kubernetes node's IP address on a statically assigned port.
# service:
    type: NodePort
    port: 80
  ingress:
    enabled: false
# $ helm install viz vizualplatform/ --values vizualplatform/values.yaml
# here viz is the release name.you can see the following output
# Release “viz” has been upgraded. Happy Helming!
# $ helm list
# you can see the release name and where the charts are deployed.



## OPERATOR-SDK:

# $ makedir helm-operator
# $ cd helm-operator
# Use the operator SDK to initialize the project. Specify the plugin and API group as the parameters for this command.initialize the helm plugin and domain nameis vizualplatform.com
# $ operator-sdk init --plugins=helm --domain vizualplatform.com
# $ operate-sdk create api
# This creates the config directory, watches.yaml and the place holder for the helm chart
# $ tree
# $ operator-sdk create api --helm-chart=vizualplatform --helm-chart-repo=Devops# $ kubectl get crd
# $ make docker-build docker-push
# $ docker images
# $ make deploy
# $ kubectl get pod -n crossplane-system
# $ kubectl get deploy -n crossplane-system




