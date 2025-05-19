# Profile

```yaml
apiVersion: skaffold/v4beta10
kind: Config
metadata:
  name: prd
requires:
  - path: ./path/to/required/skaffold.yaml
    configs: [cfg1, cfg2]
    activeProfiles:
      - name: profile1
        activatedBy: [prd]
```
