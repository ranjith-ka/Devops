# Crossplane with Azure

Create the Kube cluster in local and run the below make command to install the crossplane.

Step 1: Install the crossplane using helm charts

```bash
$ make kind crossplane
```

Step 2: Export the Subscription ID and create own subscription, this might be one time.

```bash
$ az login
$ export SUB=XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXXXXX
$ az ad sp create-for-rbac --name ranjith --sdk-auth --role Owner --scopes="/subscriptions/$SUB" -n "ranjith" > "creds.json"
$ export AZURE_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
$ az ad app permission admin-consent --id "${AZURE_CLIENT_ID}"
```

Step 3:

1.  Create Secrets and Azure Provider.

```bash
$ kubectl apply -f minikube/crossplane/azure-provider.yaml  ## Wait for the resource to be created.
$ cp minikube/crossplane/provider_example.yaml minikube/crossplane/provider.yaml
$ base64 creds.json | tr -d "\n" ## (replace in provider secrets)
$ kubectl apply -f minikube/crossplane/provider.yaml
$ kubectl apply -f minikube/crossplane/test-rg.yaml
```

Check in portal for the RG

`https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/ApplicationsListBlade`
`https://doc.crds.dev/github.com/crossplane/provider-azure@v0.19.0` ## For futher resource creation

### Install crossplane plugin

```bash
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
# Move the crossplane kubectl extension to the bin
sudo mv kubectl-crossplane /usr/local/bin
# verify that it is installed
kubectl crossplane --help
```

```bash
$ k get pkg
NAME                                                            INSTALLED   HEALTHY   PACKAGE                                                    AGE
configuration.pkg.crossplane.io/xp-getting-started-with-azure   True        True      registry.upbound.io/xp/getting-started-with-azure:latest   11m

NAME                                                   INSTALLED   HEALTHY   PACKAGE                             AGE
provider.pkg.crossplane.io/crossplane-provider-azure   True        True      crossplane/provider-azure:v0.19.0   11m
```

`az ad sp create-for-rbac --sdk-auth --role Owner --scopes="/subscriptions/XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXXXXX" -n "crossplane-sp-rbac" > "creds.json"`

`az ad sp create-for-rbac -n "MyApp"` // Custom name for the App login

```bash
$ az ad sp create-for-rbac -n "MyApp" --role Contributor --scopes /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup1} /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup2}`  ## Scope to the RG
```

Example:

`az ad sp create-for-rbac --name ranjith --role Owner --scopes="/subscriptions/$SUB" -n "ranjith" > "creds.json"`

Note: "ranjith" is the custom name for the RBAC

### References

<https://www.techtarget.com/searchitoperations/tutorial/Step-by-step-guide-to-working-with-Crossplane-and-Kubernetes>
https://crossplane.io/docs/v1.8/cloud-providers/azure/azure-provider.html
https://doc.crds.dev/github.com/crossplane/provider-azure@v0.19.0
<https://freshbrewed.science/2021/10/20/crossplane.html>
