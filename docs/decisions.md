# Decisions

## ADR-001: Use Operator Pattern over Helm Charts
**Date:** 2026-08-29  
**Status:** Accepted

**Context:**  
We need a way to provision Kubernetes environments (namespaces, RBAC, quotas) consistently for developers. Helm templates can do this once during deployment but cannot prevent configuration drift if cluster admins or users modify the resources afterward.

**Decision:**  
We will use the Kubernetes Operator Pattern built with `kubebuilder`/`controller-runtime`. A custom CRD `PlatformService` will be the single interface for developers.

**Consequences:**  
- ✅ Positive outcome: Self-healing architecture that immediately reverts any manual tampering of critical resources like NetworkPolicies.
- ✅ Positive outcome: Easy for developers—only 5 lines of YAML required.
- ⚠️ Trade-off: Higher maintenance cost than a Helm chart, requiring Go knowledge to update logic.

## ADR-002: Use unstructured.Unstructured for Dynamic Reconciliation
**Date:** 2026-08-29  
**Status:** Accepted

**Context:**  
The operator logic originally lacked typed Go structs for the `GoldenPath` CRD, causing a scaffolding failure. Rather than regenerating the entire operator schema manually for a demo repository, we need a robust method to watch and process CRs.

**Decision:**  
We will utilize `k8s.io/apimachinery/pkg/apis/meta/v1/unstructured` to dynamically watch the `GoldenPath` CRD by GVK (Group, Version, Kind) and parse its properties.

**Consequences:**  
- ✅ Positive outcome: Circumvents the need for deep code generation (`deepcopy`, `clientset`) in a lightweight demo.
- ⚠️ Trade-off: Loss of static type safety for the `GoldenPath` object fields during compilation.

## ADR-003: Declarative Multi-Cluster GitOps Packaging with ArgoCD and Kustomize
- **Date:** 2026-10-27
- **Status:** Accepted
- **Context:** Deploying this service across multi-cluster environments (`dev`, `staging`, `prod`) requires a standardized cloud-native packaging format that separates environment-agnostic base specifications from environment-specific resource allocations, replica counts, and security policies without modifying source code or inducing configuration drift.
- **Decision:** Adopt Kustomize overlays under `deploy/k8s/` and declarative ArgoCD `Application` definitions under `deploy/argocd/`.
  1. `deploy/k8s/base`: Declares non-root PodSecurity standards (`runAsUser: 10001`, `drop: [ALL]`, `readOnlyRootFilesystem: true`), probes, and resource constraints without hardcoding namespace.
  2. `deploy/k8s/overlays/{dev,staging,prod}`: Overrides namespace, replica counts, and resources per environment tier.
  3. `deploy/argocd/`: Standalone `Application` CRDs with automated self-healing in `dev`/`staging` and gated PR synchronization in `prod`.
- **Consequences:**
  - *Pros:* Fully declarative, CNCF-compliant, drift-resistant, and natively integrates with the `multi-cluster-gitops-promotion` catalog.
  - *Cons:* Requires maintaining separate overlay directories for each target cluster tier.
