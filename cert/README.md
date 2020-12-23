# Cert-Manager

To install the cert manager via Helm, set true to CustomResource/CRD's

```bash
helm install cert-manager  jetstack/cert-manager --set installCRDs=true
```

or fetch the chart and update the values.yaml

```bash
helm fetch jetstack/cert-manager
```

## Reference

<https://medium.com/flant-com/cert-manager-lets-encrypt-ssl-certs-for-kubernetes-7642e463bbce>
<https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-with-cert-manager-on-digitalocean-kubernetes>
