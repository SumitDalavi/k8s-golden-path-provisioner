# Runbook — k8s-golden-path-provisioner
> Last updated: 2026-08-29

## Prerequisites
| Tool | Required Version | How to check |
|---|---|---|
| Go | >= 1.21 | `go version` |
| kubectl | >= 1.28 | `kubectl version --client` |
| kind | Latest | `kind version` |
| Docker | >= 24.x | `docker version` |

## Quick Start
```bash
# Start a cluster
kind create cluster --name golden-path

# Apply CRD
kubectl apply -f api/v1alpha1/

# Run operator
go run . &
```

## Run Tests
```bash
# Unit tests
go test ./... -v

# E2E Tests
bash e2e/run_tests.sh
```

Expected output:
```
?       github.com/SumitDalavi/k8s-golden-path-provisioner      [no test files]
ok      github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1 0.011s
ok      github.com/SumitDalavi/k8s-golden-path-provisioner/controllers  0.252s
```

## Environment Variables
| Variable | Default | Purpose |
|---|---|---|
| KUBECONFIG | `~/.kube/config` | Kubernetes cluster credentials |
| METRICS_ADDR | `:8080` | Port for prometheus metrics |

## Common Failure Modes
| Symptom | Cause | Fix |
|---|---|---|
| `connection refused` on `go run .` | KUBECONFIG not set or cluster not running | Ensure `kind` cluster is active |
| `no matches for kind "PlatformService"` | CRD not applied | Run `kubectl apply -f api/v1alpha1/` |
