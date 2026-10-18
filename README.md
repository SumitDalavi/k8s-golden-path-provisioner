# Self-Service Kubernetes "Golden Path" Operator

[![CI Pipeline](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/ci.yml)

> **Maturity:** Full Prototype  
> A Kubernetes operator that provisions a baseline namespace, RBAC binding, and resource quota from a `PlatformService` CRD.

## The problem

Creating a new service environment usually requires several Kubernetes resources: a namespace, RBAC for deployment access, and guardrails to prevent noisy neighbours. Doing this manually is repetitive and easy to get wrong.

## The solution

`PlatformService` is a custom resource that declares the service name, owning team, and tier. The operator reconciles it into a baseline environment:

1. Creates a dedicated namespace named `<service-name>-prod`.
2. Creates a `cicd-deployer-binding` RoleBinding for the namespace's default ServiceAccount.
3. Creates a `default-quota` ResourceQuota with CPU and memory limits.

```yaml
apiVersion: platform.example.com/v1alpha1
kind: PlatformService
metadata:
  name: payments-api
spec:
  team: checkout
  tier: backend
  exposeExternally: false
```

## Current scope

| Capability | Status | Details |
|---|---|---|
| Namespace provisioning | Implemented | Creates `<name>-prod` |
| CI/CD RBAC binding | Implemented | Binds the namespace default ServiceAccount to the `edit` ClusterRole |
| Resource quota | Implemented | Sets CPU and memory limits |
| NetworkPolicies | Not implemented | Planned; the current controller does not create them |
| LimitRanges | Not implemented | Planned |
| Dedicated ServiceAccounts | Not implemented | Planned |
| Ingress / Gateway | Not implemented | `exposeExternally` is accepted by the CRD but is not reconciled |
| Tier-specific policy | Not implemented | `tier` is validated but does not change reconciliation behaviour |

## Why an operator

Helm charts and Terraform are useful for deployment-time provisioning. An operator is useful when the desired state should be continuously reconciled: if a managed resource is deleted, the controller can recreate it on the next reconciliation.

This repository demonstrates that pattern with a deliberately small, verifiable scope.

## Tech stack

- Go
- Kubebuilder / controller-runtime
- Kubernetes CRDs
- GitHub Actions

## Decision log

| Decision | Rationale |
|---|---|
| Operator over Helm | Helm is primarily a templating and release tool. An operator continuously reconciles desired state. |
| CRD abstraction | Developers declare intent through `PlatformService`; the controller creates the supporting resources. |
| Small initial scope | The first version focuses on namespace, RBAC, and quota provisioning so each behaviour can be verified. |

## Project structure

```text
api/v1alpha1/
  platformservice_types.go
controllers/
  platformservice_controller.go
config/crd/
  platformservice.yaml
docs/
README.md
```

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | >= 1.21 | Build and run the operator |
| kubectl | >= 1.28 | Kubernetes CLI |
| kind or minikube | Latest | Local Kubernetes cluster |
| Docker | >= 24.x | Container runtime |

## Run locally

```bash
git clone https://github.com/SumitDalavi/k8s-golden-path-provisioner.git
cd k8s-golden-path-provisioner

kind create cluster --name golden-path
kubectl apply -f config/crd/platformservice.yaml

go run . &
```

Create a `PlatformService`:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: platform.example.com/v1alpha1
kind: PlatformService
metadata:
  name: payments-api
spec:
  team: checkout
  tier: backend
  exposeExternally: false
EOF
```

## Verify

```bash
# CRD is registered
kubectl get crds | grep platformservice

# CR is accepted
kubectl get platformservices

# Namespace is created
kubectl get namespace payments-api-prod

# RoleBinding and ResourceQuota are created
kubectl get rolebinding,resourcequota -n payments-api-prod
```

Expected result:

```text
NAME                  STATUS
payments-api-prod     Active

NAME                                              ROLE                AGE
cicd-deployer-binding                             ClusterRole/edit    ...

NAME             AGE
default-quota    ...
```

## Cleanup

```bash
kind delete cluster --name golden-path
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Runbook](docs/runbook.md)
- [Decisions](docs/decisions.md)
- [Changelog](docs/changelog.md)

## Author

**Sumit Dalavi** — Senior Platform / DevSecOps Engineer  
[GitHub](https://github.com/SumitDalavi) | [LinkedIn](https://in.linkedin.com/in/sumit-dalavi-762838129)

---

Built to demonstrate continuous reconciliation patterns with an honest, testable scope.
