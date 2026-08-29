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
