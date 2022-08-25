# Crossplane with Azure

Create the Kube cluster in local and run the below make command to install the crossplane.

Step1: Install the crossplane using helm charts

```bash
make crossplane
```

Install crossplane plugin

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

### References

    - https://www.techtarget.com/searchitoperations/tutorial/Step-by-step-guide-to-working-with-Crossplane-and-Kubernetes
    - https://crossplane.io/docs/v1.8/cloud-providers/azure/azure-provider.html
