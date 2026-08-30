#!/usr/bin/env bash
set -euo pipefail

KIND_CLUSTER="provisioner-e2e"
log() { echo "[e2e] $*"; }

log "Creating KIND cluster: $KIND_CLUSTER"
kind create cluster --name "$KIND_CLUSTER" --wait 60s

log "Installing CRD..."
kubectl apply -f config/crd/

log "Starting controller in background..."
go run main.go &
PID=$!
sleep 5 # give it a moment to start and connect

log "Applying a GoldenPath CR..."
kubectl apply -f e2e/fixtures/goldenpath-sample.yaml

log "Waiting for reconciliation..."
sleep 10

# Assert: resources were provisioned
log "Asserting resources..."
kubectl get namespace sample || { log "❌ Namespace 'sample' not created"; kill $PID; exit 1; }
kubectl get serviceaccount -n sample default || { log "❌ SA not created"; kill $PID; exit 1; }

log "✅ Provisioner reconciled GoldenPath CR successfully"
kill $PID
kind delete cluster --name "$KIND_CLUSTER"
