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
