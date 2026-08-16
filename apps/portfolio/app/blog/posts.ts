export type PostSection = { heading: string; body: string; code?: string; note?: string };
export type BlogPost = { slug: string; title: string; excerpt: string; category: "KIND" | "Minikube"; level: string; readTime: string; date: string; source: string; accent: string; sections: PostSection[] };

export const posts: BlogPost[] = [
  {
    slug: "kubernetes-in-docker-with-kind", title: "Build a local Kubernetes platform with KIND", category: "KIND", level: "Beginner", readTime: "6 min read", date: "Aug 2026", accent: "lime",
    excerpt: "Create a repeatable multi-node Kubernetes environment, connect a local registry and route traffic through ingress.",
    source: "https://github.com/ranjith-ka/Devops/blob/main/kind/README.md",
    sections: [
      { heading: "Why KIND?", body: "KIND runs Kubernetes nodes as Docker containers. It is fast, disposable and ideal for testing platform configuration, Helm charts and delivery workflows before using a shared cluster." },
      { heading: "Create the cluster", body: "With Docker running, install KIND and Helm, then create a named cluster using the repository configuration.", code: "brew install kind helm\nkind create cluster --config kind/config.yaml --name k8s\nkubectl cluster-info\nkubectl get nodes" },
      { heading: "Use a local registry", body: "Connect the registry container to KIND's Docker network, tag the image and push it to localhost:5000.", code: "docker network connect kind kind-registry\ndocker tag myimage:tag localhost:5000/myimage:tag\ndocker push localhost:5000/myimage:tag\ncurl http://localhost:5000/v2/_catalog", note: "The KIND nodes must be configured to recognize the local registry. The repository includes working registry and cluster examples." },
      { heading: "Install the sample platform", body: "The Makefile turns the setup into a repeatable platform workflow.", code: "make snapshot\nmake kind-cluster\nmake load-image\nmake ingress\nmake install-app" },
      { heading: "Verify canary routing", body: "Send different request headers through ingress to verify production and canary workloads.", code: "curl -s -H \"testing: always\" http://awesome-http.example.com/dev\ncurl -s -H \"testing: never\" http://awesome-http.example.com/dev" },
    ],
  },
  {
    slug: "gateway-api-on-kind", title: "Gateway API on KIND: from Ingress to HTTPRoute", category: "KIND", level: "Intermediate", readTime: "8 min read", date: "Aug 2026", accent: "orange",
    excerpt: "Install Gateway API CRDs and Envoy Gateway, then expose a local application using Gateway and HTTPRoute resources.",
    source: "https://github.com/ranjith-ka/Devops/blob/main/kind/GATEWAY_API_SETUP.md",
    sections: [
      { heading: "What changes from Ingress?", body: "Gateway API separates infrastructure ownership from application routing. GatewayClass selects the controller, Gateway defines listeners and HTTPRoute attaches application rules." },
      { heading: "Prepare KIND", body: "Recreate the cluster with the feature gates and port mappings defined in the repository configuration.", code: "kind delete cluster --name k8s\nkind create cluster --config kind/config.yaml --name k8s\nkubectl get nodes" },
      { heading: "Install the APIs", body: "Apply the standard Gateway API CRDs and confirm the resource definitions.", code: "kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml\nkubectl get crds | grep gateway" },
      { heading: "Install Envoy Gateway", body: "A controller is required to implement the API and create the data plane.", code: "helm repo add envoy-gateway https://gateway.envoyproxy.io\nhelm repo update\nhelm install envoy-gateway envoy-gateway/gateway-helm --namespace envoy-gateway-system --create-namespace\nkubectl get gatewayclasses" },
      { heading: "Test a route", body: "After applying Gateway and HTTPRoute resources, port-forward and send the configured Host header.", code: "kubectl get gateway my-gateway\nkubectl port-forward svc/my-gateway 8080:80\ncurl -H \"Host: awesome-http.example.com\" http://localhost:8080/dev" },
    ],
  },
  {
    slug: "spiffe-spire-workload-identity-kind", title: "Workload identity with SPIFFE and SPIRE on KIND", category: "KIND", level: "Advanced", readTime: "12 min read", date: "Aug 2026", accent: "mint",
    excerpt: "Understand SPIFFE IDs and SVIDs, deploy SPIRE Server and Agents, and issue short-lived identities to Kubernetes workloads.",
    source: "https://github.com/ranjith-ka/Devops/blob/main/kind/SPIFFE_SPIRE_TUTORIAL.md",
    sections: [
      { heading: "Identity instead of location", body: "SPIFFE gives a workload a stable URI identity that survives IP, node and pod changes. An SVID proves the identity while the Workload API delivers it locally." },
      { heading: "Create the lab", body: "Start the KIND cluster and clone the official SPIRE tutorial manifests.", code: "kind create cluster --config kind/config.yaml --name k8s\ncd /tmp\ngit clone --depth 1 https://github.com/spiffe/spire-tutorials.git\ncd spire-tutorials/k8s/quickstart" },
      { heading: "Deploy SPIRE Server", body: "The server acts as the certificate authority and policy authority for the trust domain.", code: "kubectl apply -f spire-namespace.yaml\nkubectl apply -f server-account.yaml -f spire-bundle-configmap.yaml -f server-cluster-role.yaml\nkubectl apply -f server-configmap.yaml -f server-statefulset.yaml -f server-service.yaml\nkubectl -n spire rollout status statefulset/spire-server" },
      { heading: "Deploy node agents", body: "SPIRE Agents run as a DaemonSet, attest workloads and expose the local Workload API socket.", code: "kubectl apply -f agent-account.yaml -f agent-cluster-role.yaml\nkubectl apply -f agent-configmap.yaml -f agent-daemonset.yaml\nkubectl -n spire get daemonset,pods" },
      { heading: "Verify identity delivery", body: "A registration entry must match the workload namespace and service account. The workload can then receive a short-lived X.509-SVID without storing a long-lived secret.", note: "The full source tutorial continues through registration entries, SVID retrieval and certificate inspection." },
    ],
  },
  {
    slug: "minikube-local-platform-setup", title: "Turn Minikube into a local platform lab", category: "Minikube", level: "Beginner", readTime: "5 min read", date: "Aug 2026", accent: "blue",
    excerpt: "Use Minikube as a lightweight Docker Desktop alternative and install ingress, applications and monitoring through one workflow.",
    source: "https://github.com/ranjith-ka/Devops/blob/main/minikube/README.md",
    sections: [
      { heading: "Build a repeatable lab", body: "The repository Makefile captures cluster creation and platform add-ons so the environment can be rebuilt whenever required.", code: "make minikube\neval $(minikube -p minikube docker-env)\nmake snapshot\nmake ingress\nmake install-app\nmake monitoring" },
      { heading: "Build directly into Minikube", body: "Point your shell at Minikube's Docker daemon so local images work without a remote registry.", code: "eval $(minikube -p minikube docker-env)\ndocker build -t my-app:dev .\nkubectl get pods --watch" },
      { heading: "Use a proxy-aware cluster", body: "Pass proxy settings into the Docker environment and include the cluster subnet in NO_PROXY.", code: "minikube start --docker-env HTTP_PROXY=http://192.168.1.9:10809 --docker-env HTTPS_PROXY=http://192.168.1.9:10809 --docker-env NO_PROXY=192.168.99.0/24\nexport no_proxy=$no_proxy,$(minikube ip)" },
      { heading: "Reset cleanly", body: "A local platform should be disposable. Delete and rebuild through the same Make targets when experiments leave unknown state.", code: "minikube delete" },
    ],
  },
  {
    slug: "flux-gitops-on-minikube", title: "Practice Flux GitOps locally with Minikube", category: "Minikube", level: "Intermediate", readTime: "7 min read", date: "Aug 2026", accent: "lime",
    excerpt: "Install Flux controllers, apply Git sources and Helm releases, then reconcile changes and inspect the deployed artifact.",
    source: "https://github.com/ranjith-ka/Devops/blob/main/minikube/flux/Readme.md",
    sections: [
      { heading: "Learn the controllers", body: "Flux is composed of source, kustomize, helm, notification and image automation controllers. Practising locally makes production reconciliation easier to reason about." },
      { heading: "Install and apply", body: "Start the local cluster, install Flux and apply the staging resources.", code: "make kind\nflux install\nkubectl apply -f minikube/flux/staging\nkubectl get pods" },
      { heading: "Reconcile immediately", body: "Flux normally polls Git. During development, request reconciliation and inspect the resulting resources.", code: "flux reconcile source git devops\nflux get all\nhelm list\nkubectl get pods -o wide" },
      { heading: "Verify the artifact", body: "A Helm chart version and a container image are different signals. Inspect the pod to confirm the deployed tag or digest.", code: "kubectl describe pod <pod-name>\nkubectl get pod <pod-name> -o jsonpath='{.spec.containers[*].image}'" },
    ],
  },
  {
    slug: "keda-event-driven-autoscaling-minikube", title: "Explore event-driven autoscaling with KEDA", category: "Minikube", level: "Intermediate", readTime: "6 min read", date: "Aug 2026", accent: "orange",
    excerpt: "Understand how KEDA extends Kubernetes autoscaling for queues, streams and event sources in a safe local environment.",
    source: "https://github.com/ranjith-ka/Devops/tree/main/minikube/keda",
    sections: [
      { heading: "Why KEDA?", body: "Kubernetes HPA works well with CPU and memory. KEDA adds event-aware scalers for Kafka, RabbitMQ and cloud queues while using standard Kubernetes primitives." },
      { heading: "The operating model", body: "A ScaledObject links a Deployment to an event trigger. KEDA monitors the trigger, activates from zero and hands scaling decisions to an HPA." },
      { heading: "Install the operator", body: "Install KEDA into its own namespace and confirm the operator and metrics server are ready.", code: "helm repo add kedacore https://kedacore.github.io/charts\nhelm repo update\nhelm install keda kedacore/keda --namespace keda --create-namespace\nkubectl get pods -n keda" },
      { heading: "Observe the loop", body: "Apply the repository example, generate events and watch the Deployment, HPA and KEDA resources together.", code: "kubectl apply -f minikube/keda/tutorial/demo.yaml\nkubectl get scaledobjects,hpa\nkubectl get deployment --watch", note: "Treat scale-to-zero carefully for latency-sensitive services. Configure activation thresholds, fallback and authentication explicitly." },
    ],
  },
];

export function getPost(slug: string) { return posts.find((post) => post.slug === slug) }
