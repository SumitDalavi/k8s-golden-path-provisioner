# PlatformService CRD Reference

## Spec

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `image` | string | Yes | — | Container image (e.g. `nginx:1.25`) |
| `replicas` | int | No | `2` | Number of pod replicas (1–20) |
| `port` | int | No | `8080` | Container port to expose |
| `env` | map[string]string | No | `{}` | Environment variables injected into pods |
| `serviceType` | string | No | `ClusterIP` | Kubernetes Service type |

## Status

| Field | Type | Description |
|-------|------|-------------|
| `phase` | string | `Provisioning`, `Running`, `Degraded`, `Failed` |
| `readyReplicas` | int | Number of ready pod replicas |
| `conditions` | []Condition | Standard K8s conditions array |
| `observedGeneration` | int64 | Last reconciled generation |

## Validation Rules
- `replicas` must be 1–20
- `port` must be 1–65535
- `image` must be non-empty

## Example

```yaml
apiVersion: platform.example.com/v1alpha1
kind: PlatformService
metadata:
  name: my-api
spec:
  image: myapp:v1.2.3
  replicas: 3
  port: 8080
  env:
    LOG_LEVEL: info
    DB_HOST: postgres.default.svc
```
