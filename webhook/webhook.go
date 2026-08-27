package webhook

import (
	"context"
	"fmt"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PlatformServiceDefaulter implements the admission.CustomDefaulter interface.
// It sets sensible defaults on PlatformService resources before creation.
type PlatformServiceDefaulter struct{}

func (d *PlatformServiceDefaulter) Default(ctx context.Context, obj client.Object) error {
	ps, ok := obj.(*platformv1.PlatformService)
	if !ok {
		return fmt.Errorf("expected PlatformService, got %T", obj)
	}
	if ps.Spec.Tier == "" {
		ps.Spec.Tier = "backend" 
	}
	if ps.Spec.Team == "" {
		ps.Spec.Team = "default-team"
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
	if ps.Spec.Tier != "frontend" && ps.Spec.Tier != "backend" && ps.Spec.Tier != "worker" {
		return fmt.Errorf("tier must be frontend, backend, or worker")
	}
	if ps.Spec.Team == "" {
		return fmt.Errorf("team must be specified")
	}
	return nil
}
