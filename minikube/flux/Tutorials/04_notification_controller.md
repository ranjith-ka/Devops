### Notification controller

1. To understand in simple words, there is provider and Alert linked to the provider for Notification controller to work. 

## Alerts

Discord webhook

- Simple to create a discord webhook
- Below webook will be removed, used only for demo purpose only
  https://discord.com/api/webhooks/1068767766733144175/NiCgkH7MiRlKAJs1f15uUyikbMCrt-3LRuKPe8yeDQcfhMr_so3Ydgzpnv8PBZ5qbvXw

```bash
kubectl -n default create secret generic discord-url \
--from-literal=address=https://discord.com/api/webhooks/1068767766733144175/NiCgkH7MiRlKAJs1f15uUyikbMCrt-3LRuKPe8yeDQcfhMr_so3Ydgzpnv8PBZ5qbvXw
```

- Test the webhook as well

```bash
curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X POST --data "{\"content\": \"Posted Via Command line\"}" https://discord.com/api/webhooks/1068767766733144175/NiCgkH7MiRlKAJs1f15uUyikbMCrt-3LRuKPe8yeDQcfhMr_so3Ydgzpnv8PBZ5qbvXw
```

```bash
$ kubectl apply --server-side -f minikube/flux/notification.yaml
```

## --server-side (did the magic for me)

### Events

### Incoming Webhook

### Notification API

### Providers
