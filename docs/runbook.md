# Runbook — k8s-golden-path-provisioner
> Last updated: 2026-10-18

## Prerequisites
| Tool | Required Version | How to check |
|---|---|---|
| Go | >= 1.26 | `go version` |
| kubectl | >= 1.28 | `kubectl version --client` |
| kind | Latest | `kind version` |
| Docker | >= 24.x | `docker version` |

## Quick Start
```bash
# Start a cluster
kind create cluster --name golden-path

# Apply CRD
kubectl apply -f config/crd/platformservice.yaml

# Run operator
go run . &
```

## Run Tests
```bash
# Unit tests
go test ./... -v

# E2E Tests
bash e2e/test_provisioner_e2e.sh
```

Expected output:
```
ok      github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1 0.011s
ok      github.com/SumitDalavi/k8s-golden-path-provisioner/webhook     0.011s
ok      github.com/SumitDalavi/k8s-golden-path-provisioner/tests       0.012s
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
| `no matches for kind "PlatformService"` | CRD not applied | Run `kubectl apply -f config/crd/platformservice.yaml` |

## Release Runbook

### Triggering a New Release
Releases are automated via GitHub Actions (`.github/workflows/release.yml`) and trigger whenever a semantic version tag is pushed:

```bash
# 1. Ensure master is clean and up-to-date
git checkout master
git pull origin master

# 2. Tag the release
git tag -a v0.2.0 -m "Release v0.2.0"

# 3. Push the tag to trigger the pipeline
git push origin v0.2.0
```

### Release Pipeline Stages
1. **Test Verification:** Runs `go test ./...` across all packages.
2. **Container Image Build:** Builds the minimal operator image (`ghcr.io/sumitdalavi/k8s-golden-path-provisioner:<tag>`).
3. **Trivy Vulnerability Scan:** Scans the built container for CRITICAL and HIGH CVEs.
4. **Anchore Syft SBOM Generation:** Creates an SPDX JSON Software Bill of Materials (`sbom.spdx.json`).
5. **GitHub Release Publishing:** Creates the release with auto-generated release notes and attaches `sbom.spdx.json`.

### Verifying the Release
```bash
# Check workflow run status
gh run list --workflow=release.yml

# Inspect the published release and attached assets
gh release view v0.1.0

# Download and inspect the SBOM
gh release download v0.1.0 -p "sbom.spdx.json"
cat sbom.spdx.json | jq .spdxVersion
```

### Release Failure Modes
| Symptom | Cause | Resolution |
|---|---|---|
| `Unable to resolve action ... unable to find version` | Pinned action tag missing or upstream tag modified | Pin action to a verified release tag or commit SHA |
| `invalid reference format: repository name must be lowercase` | Container repository owner contains uppercase characters | Lowercase the repository owner before tagging |
| `File not found: sbom.spdx.json` during release step | SBOM action configured with `artifact-name` instead of `output-file` | Ensure `output-file: sbom.spdx.json` is set in Anchore action |

## Operations & Diagnostics Guide (v0.2.0)

### Mode Verification & Diagnostics
Inspect the active execution mode and backing provider status via the health endpoint:
```bash
# Check runtime mode and adapter status
curl -s http://localhost:8080/health | jq .
# Expected output when running without keys:
# {"status": "ok", "mode": {"database": "mock", "ai": "mock", "cache": "mock"}}
```

### Running in Mock Mode (Tier 1)
Start the service with zero external dependencies:
```bash
cp .env.example .env
# Start service without modifying .env
npm start || python main.py || go run .
```

### Running with Live Credentials (Tier 2 / 3)
Populate target credentials in `.env` and restart:
```bash
# Verify live upgrade via health check
curl -s http://localhost:8080/health | jq .mode
```

### Pulling Pre-Built OCI Image
```bash
docker pull ghcr.io/sumitdalavi/k8s-golden-path-provisioner:v0.2.0
docker run -d -p 8080:8080 --env-file .env ghcr.io/sumitdalavi/k8s-golden-path-provisioner:v0.2.0
```

## GitOps Deployment & Multi-Cluster Operations (v0.2.1)

### Local Overlay Validation
Validate environment-specific manifests offline with `kubectl kustomize`:
```bash
# Validate Dev Overlay
kubectl kustomize deploy/k8s/overlays/dev

# Validate Staging Overlay
kubectl kustomize deploy/k8s/overlays/staging

# Validate Production Overlay
kubectl kustomize deploy/k8s/overlays/prod
```

### ArgoCD Standalone Application Deployment
Apply the standalone declarative ArgoCD application to your local or target cluster:
```bash
# Deploy default staging application
kubectl apply -f deploy/argocd/application.yaml

# Deploy environment-specific applications
kubectl apply -f deploy/argocd/application-dev.yaml
kubectl apply -f deploy/argocd/application-prod.yaml
```
