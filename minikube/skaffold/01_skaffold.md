# Skaffold

<https://skaffold.dev/docs/>

Skaffold is a tool that facilitates continuous development for Kubernetes applications.
It can be used to build, test, and deploy applications.

## Skaffold Workflow and Architecture

![alt text](https://skaffold.dev/images/architecture.png)

### Install

```bash
brew install skaffold
```

## skaffold dev

skaffold dev enables continuous local development on an application. While in dev mode, Skaffold will watch an application’s source files, and when it detects changes, will rebuild your images (or sync files to your running containers), push any new images, test built images, and redeploy the application to your cluster.

Dev Loop:

- File sync
- build
- Test
- Deploy

Skaffold also supports a polling mode where the filesystem is checked for changes on a configurable interval, or a manual mode, where Skaffold waits for user input to check for file changes. These watch modes can be configured through the --trigger flag.

## Debugging With Skaffold

### TODO

- [ ] Add debugging with Skaffold
- [ ] Steps for Golang and .NET Core

## Continuous Delivery

Skaffold can be used to automate the deployment of applications to Kubernetes. Skaffold can be used to deploy applications to a Kubernetes cluster, and can also be used to deploy applications to a local Kubernetes cluster, such as Minikube.

### Deploying to a Kubernetes Cluster

To deploy an application to a Kubernetes cluster, you can use the skaffold deploy command. This command will build the application, push the image to a registry, and deploy the application to the Kubernetes cluster.

```bash
skaffold deploy
```

