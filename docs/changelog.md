# Changelog

## [v0.1.0] — 2026-10-18
### Added
- Automated GitHub Actions release pipeline triggered on `v*.*.*` git tags (`.github/workflows/release.yml`).
- Container vulnerability scanning step via Aqua Security Trivy (`aquasecurity/trivy-action@v0.36.0`).
- SPDX-format Software Bill of Materials (SBOM) generation step (`anchore/sbom-action@v0`).
- Automated GitHub release publishing attaching `sbom.spdx.json` via `softprops/action-gh-release@v2`.

### Fixed
- Upgraded `aquasecurity/trivy-action` to `v0.36.0` to eliminate broken upstream `setup-trivy@v0.2.1` tag reference failures.
- Sanitized container image repository tag to ensure all-lowercase casing compliance (`SumitDalavi` -> `sumitdalavi`) per OCI image specification.
- Configured `output-file: sbom.spdx.json` in `anchore/sbom-action` so the generated SBOM artifact is written to disk for release asset upload.
- Aligned Go version to `1.26` across release pipeline and project runtime dependencies.

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

### Post-Release Hotfixes
- Resolved lingering CI/CD failures introduced during portfolio elevation.
- Fixed Docker build and permission errors across client/server components.
- Corrected Kubernetes controller GroupVersionKind mismatches and E2E Vault addressing.
- Repaired broken property-based test configurations.
