# Kuma Deployments

    ```bash
    $ helm install --namespace kuma-system --set "controlPlane.mode=zone" kuma kuma/kuma
    ```

## Single-zone deployment

### Deploy a single-zone control plane