package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// PlatformServiceDefaulter implements the admission.CustomDefaulter interface.
// It sets sensible defaults on PlatformService resources before creation.
type PlatformServiceDefaulter struct{}

func (d *PlatformServiceDefaulter) Default(ctx context.Context, obj client.Object) error {
	ps, ok := obj.(*platformv1.PlatformService)
	if !ok {
		return fmt.Errorf("expected PlatformService, got %T", obj)
	}
	if ps.Spec.Replicas == 0 {
		ps.Spec.Replicas = 2 // default to 2 replicas for HA
	}
	if ps.Spec.Port == 0 {
		ps.Spec.Port = 8080
	}
	if ps.Spec.Image == "" {
		ps.Spec.Image = "nginx:1.25-alpine"
	}
	if ps.Labels == nil {
		ps.Labels = make(map[string]string)
	}
	ps.Labels["app.kubernetes.io/managed-by"] = "golden-path-provisioner"
	return nil
}

// PlatformServiceValidator implements admission.CustomValidator.
// It enforces invariants that the CRD schema alone cannot express.
type PlatformServiceValidator struct{}

func (v *PlatformServiceValidator) ValidateCreate(ctx context.Context, obj client.Object) error {
	return v.validate(obj)
}

func (v *PlatformServiceValidator) ValidateUpdate(ctx context.Context, oldObj, newObj client.Object) error {
	return v.validate(newObj)
}

func (v *PlatformServiceValidator) ValidateDelete(ctx context.Context, obj client.Object) error {
	return nil // deletion is always allowed
}

func (v *PlatformServiceValidator) validate(obj client.Object) error {
	ps, ok := obj.(*platformv1.PlatformService)
	if !ok {
		return fmt.Errorf("expected PlatformService, got %T", obj)
	}
	if ps.Spec.Replicas > 20 {
		return fmt.Errorf("replicas must be <= 20, got %d", ps.Spec.Replicas)
	}
	if ps.Spec.Replicas < 1 {
		return fmt.Errorf("replicas must be >= 1")
	}
	if ps.Spec.Port < 1 || ps.Spec.Port > 65535 {
		return fmt.Errorf("port must be 1-65535, got %d", ps.Spec.Port)
	}
	if ps.Spec.Image == "" {
		return fmt.Errorf("image must be specified")
	}
	return nil
}
