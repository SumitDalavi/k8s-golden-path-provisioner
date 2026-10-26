# Changelog

## [0.2.1] - 2026-10-27
### Added
- **Declarative GitOps Deployment:** Standardized `deploy/k8s/` base and environment overlays (`dev`, `staging`, `prod`) adhering to CNCF restricted pod security standards (non-root UID 10001, dropped capabilities, read-only rootfs).
- **ArgoCD Continuous Delivery:** Added declarative ArgoCD application manifests (`deploy/argocd/application.yaml`, `application-dev.yaml`, `application-prod.yaml`, `values-override.yaml`).
- **Multi-Cluster Orchestration:** Integrated with `multi-cluster-gitops-promotion` catalog for automated zero-drift self-healing in development/staging and gated PR promotion in production.

## [0.2.0] - 2026-10-20
### Added
- **Progressive Credential Loading (PCL):** Dynamic auto-promotion between zero-dependency mock fallbacks and live production backends based on environment credentials.
- **Dual Distribution:** Universal OCI container publishing to GitHub Container Registry (`ghcr.io/sumitdalavi/k8s-golden-path-provisioner`) and Helm OCI registries alongside GitHub Release tarballs.
- **Health Reporting:** Exposed `/health` or `/healthz` endpoint with PCL status detection (`mock` vs `live`).
- **Dynamic Reloading:** Support for atomic ConfigMap / secret reloads without requiring pod restarts.
- **Supply Chain Security:** Trivy vulnerability scan gates and SPDX JSON SBOM (`sbom.spdx.json`) generation for all releases.

## [0.1.0] - 2026-10-18
### Added
- Initial baseline service implementation, unit test suites, and CI workflows.

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
