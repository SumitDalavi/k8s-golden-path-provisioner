# Changelog

## [2026-08-29] — Phase 2 Evidence
### Added
- Created `e2e/run_tests.sh` script to verify E2E CR reconciliation and RBAC testing on `kind`.
- Checked in `e2e/results.log` as a test verification artifact.
- Updated `ARCHITECTURE.md` with complete Mermaid diagram and component/dependency tables.
- Standardized documentation (`runbook.md`, `decisions.md`).

### Fixed
- Fixed Helm templating errors (replaced undefined `include` directives with static names).
- Added `rbac.yaml` with `ServiceAccount`, `ClusterRole`, and `ClusterRoleBinding` to the Helm chart.
- Corrected `--leader-elect` flag to `--enable-leader-election` to prevent `CrashLoopBackOff`.
- Removed undefined health probes (`/healthz` and `/readyz`) from deployment to resolve restart loops.
- Rewrote Go controller logic in `main.go` and `platformservice_controller.go` to use `unstructured.Unstructured` to dynamically watch `GoldenPath` CRD and execute namespace and RBAC reconciliation.
