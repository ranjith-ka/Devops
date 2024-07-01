### Authentication

Often a scaler will require authentication or secrets and config to check for events.

KEDA provides a few secure patterns to manage authentication flows:

### Defining secrets and config maps on ScaledObject 

If using the RabbitMQ scaler, the host parameter may include passwords so is required to be a reference. 

You can create a secret with the value of the host string, reference that secret in the deployment, and map it to the ScaledObject metadata parameter like below:

TriggerAuthentication resource to define authentication as a separate resource to a ScaledObject

```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: {trigger-authentication-name}
  namespace: default # must be same namespace as the ScaledObject
spec:
  podIdentity:
      provider: none | azure | azure-workload | aws-eks | aws-kiam | gcp  # Optional. Default: none
      identityId: <identity-id>                                           # Optional. Only used by azure & azure-workload providers.
  secretTargetRef:                                                        # Optional.
  - parameter: {scaledObject-parameter-name}                              # Required.
    name: {secret-name}                                                   # Required.
    key: {secret-key-name}                                                # Required.
  env:                                                                    # Optional.
  - parameter: {scaledObject-parameter-name}                              # Required.
    name: {env-name}                                                      # Required.
    containerName: {container-name}                                       # Optional. Default: scaleTargetRef.envSourceContainerName of ScaledObject
  hashiCorpVault:                                                         # Optional.
    address: {hashicorp-vault-address}                                    # Required.
    namespace: {hashicorp-vault-namespace}                                # Optional. Default is root namespace. Useful for Vault Enterprise
    authentication: token | kubernetes                                    # Required.
    role: {hashicorp-vault-role}                                          # Optional.
    mount: {hashicorp-vault-mount}                                        # Optional.
    credential:                                                           # Optional.
      token: {hashicorp-vault-token}                                      # Optional.
      serviceAccount: {path-to-service-account-file}                      # Optional.
    secrets:                                                              # Required.
    - parameter: {scaledObject-parameter-name}                            # Required.
      key: {hasicorp-vault-secret-key-name}                               # Required.
      path: {hasicorp-vault-secret-path}  
```

If creating a deployment yaml that references a secret, be sure the secret is created before the deployment that references it, and the scaledObject after both of them to avoid invalid references.