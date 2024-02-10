# Kuma Deployments

    ```bash
    $ helm install --namespace kuma-system --set "controlPlane.mode=zone" kuma kuma/kuma
    ```

## Single-zone deployment

### Deploy a single-zone control plane

Default helm chart will install the single-zone control plane.

    ```bash
    $ helm install --create-namespace --namespace kuma-system \
    --set "controlPlane.mode=zone" \
    kuma kuma/kuma

    $ make kuma
    ```

