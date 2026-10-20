# Self-Service Kubernetes "Golden Path" Operator

[![CI Pipeline](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/SumitDalavi/k8s-golden-path-provisioner?color=0969da&label=release&logo=github)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/releases)
[![Release Workflow](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/release.yml/badge.svg)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/actions/workflows/release.yml)
[![GHCR Container](https://img.shields.io/badge/ghcr.io-v0.2.0-blue?logo=docker&logoColor=white)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/pkgs/container/k8s-golden-path-provisioner)
[![Helm OCI](https://img.shields.io/badge/helm%20chart-0.2.0-0f1689?logo=helm&logoColor=white)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/pkgs/container/k8s-golden-path-provisioner)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![SBOM](https://img.shields.io/badge/SBOM-SPDX%202.3-blueviolet?logo=spdx)](https://github.com/SumitDalavi/k8s-golden-path-provisioner/releases/tag/v0.2.0)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

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
| Go | >= 1.26 | Build and run the operator |
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

## Release & Security Verification

Every tagged release (`v*.*.*`) triggers an automated DevSecOps release workflow:
1. **Automated Testing:** Runs `go test ./...` unit tests.
2. **Container Image Build:** Builds the minimal distroless operator image (`ghcr.io/sumitdalavi/k8s-golden-path-provisioner:<tag>`).
3. **Vulnerability Scanning:** Scans the operator container image using [Aqua Security Trivy](https://github.com/aquasecurity/trivy-action) for `HIGH` and `CRITICAL` CVEs.
4. **SBOM Generation:** Generates a Software Bill of Materials (SBOM) in SPDX JSON format (`sbom.spdx.json`) using [Anchore Syft](https://github.com/anchore/sbom-action).
5. **GitHub Release Publishing:** Packages release notes and attaches `sbom.spdx.json` as a verifiable supply-chain asset.

## 🚀 Progressive Credential Loading & Execution Tiers

This project implements **Progressive Credential Loading (PCL)**, allowing it to run immediately out-of-the-box without external dependencies while seamlessly upgrading as credentials are provided.

### Tier 1: Zero-Config Mock Mode (Default)
Ideal for CI, testing, and rapid local evaluation.
- **Execution:** Run with default configuration (`cp .env.example .env`). Zero API keys, databases, or cloud credentials required.
- **Fallback Behavior:** In-memory caches, deterministic offline heuristics, mock vector embeddings, and local mock storage are automatically used.
- **Health Verification:** Endpoints report `{"mode": "mock"}` for unconfigured dependencies.

### Tier 2: Local Development with Live Services
Connect specific external services while keeping others mocked.
- **Execution:** Provide credentials in `.env` for the services you want active (e.g., `KUBECONFIG`).
- **Progressive Upgrade:** The runtime inspects environment variables on boot and dynamically promotes only configured modules to live mode.
- **Safe Fallback:** Any missing credentials automatically remain on mock fallbacks with no crashes or missing service errors.

### Tier 3: Production Hardened
Enterprise deployment with production-grade backends, security scans, and container distributions.
- **Helm OCI Chart:** Deploy directly via Helm without adding a chart repository:
  ```bash
  helm install k8s-golden-path-provisioner oci://ghcr.io/sumitdalavi/k8s-golden-path-provisioner --version 0.2.0
  ```
- **Kubernetes:** Injects secrets directly via Helm values or External Secrets Operator.
- **Release Verification:** Signed artifacts, checksums, and SPDX SBOM available on [GitHub Releases](https://github.com/SumitDalavi/k8s-golden-path-provisioner/releases).

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
