# Self-Service Kubernetes "Golden Path" Operator 🌟🛣️

> A custom Kubernetes Operator that provisions compliant, production-ready environments (Namespaces, RBAC, NetworkPolicies, ResourceQuotas) from a single Developer CRD.

## The Problem

When developers need to deploy a new service, they usually need a namespace, service accounts, network policies (to restrict blast radius), resource limits, and secrets access. Creating all of these YAML files manually is tedious, error-prone, and violates Platform Engineering principles of self-service.

## The Solution

This project implements a Custom Resource Definition (CRD) called `PlatformService`. A developer submits one simple YAML file declaring their service's name and tier (e.g., `tier: backend`). 

The Operator automatically reconciles this into a "Golden Path" environment:
1. Provisions a dedicated **Namespace**
2. Applies default **NetworkPolicies** (e.g., deny-all ingress by default)
3. Configures **ResourceQuotas** and **LimitRanges** to prevent noisy neighbors
4. Creates dedicated **ServiceAccounts** and **RoleBindings** for CI/CD deployment access

```yaml
apiVersion: platform.internal/v1alpha1
kind: PlatformService
metadata:
  name: payments-api
spec:
  team: checkout
  tier: backend
  exposeExternally: false
```

## Why This Over the Obvious Alternative

Many teams try to solve this with Helm charts or raw Terraform. While those work for *deploying* applications, they are poor abstractions for *managing cluster multi-tenancy state*. By using a Kubernetes Operator, we ensure that if someone manually deletes a NetworkPolicy, the controller instantly recreates it — enforcing compliance continuously, not just at deploy-time.

## 🛠️ Tech Stack

- **Language**: Go
- **Framework**: kubebuilder / controller-runtime
- **API**: Kubernetes Custom Resource Definitions (CRDs)

## Decision Log

| Decision | Rationale |
|----------|-----------|
| Go Operator over Helm | Helm is a one-time templating engine. An Operator runs a continuous reconciliation loop, ensuring the environment *stays* compliant even if drifted. |
| CRD Abstraction | The `PlatformService` CRD abstracts away 100s of lines of K8s boilerplate into 5 lines of declarative developer intent. |
| Default NetworkPolicies | Implementing a Zero-Trust network architecture is too complex to ask developers to do manually; the platform must enforce it automatically. |

## 📁 Project Structure

```
├── api/v1alpha1/
│   └── platformservice_types.go       # CRD definition
├── controllers/
│   └── platformservice_controller.go  # Reconciliation logic
├── docs/ARCHITECTURE.md
└── README.md
```

## 👨‍💻 Author

*Built to demonstrate Platform Engineering abstractions and Kubernetes API extension.*
