package webhook_test

import (
	"context"
	"testing"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	"github.com/SumitDalavi/k8s-golden-path-provisioner/webhook"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newPS(team, tier string) *platformv1.PlatformService {
	return &platformv1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec:       platformv1.PlatformServiceSpec{Team: team, Tier: tier},
	}
}

func TestDefaulter_SetsDefaults(t *testing.T) {
	d := &webhook.PlatformServiceDefaulter{}
	ps := newPS("", "")
	if err := d.Default(context.Background(), ps); err != nil {
		t.Fatal(err)
	}
	if ps.Spec.Tier != "backend" {
		t.Errorf("expected tier=backend, got %s", ps.Spec.Tier)
	}
	if ps.Spec.Team != "default-team" {
		t.Errorf("expected team=default-team, got %s", ps.Spec.Team)
	}
	if ps.Labels["app.kubernetes.io/managed-by"] != "golden-path-provisioner" {
		t.Error("missing managed-by label")
	}
}

func TestValidator_RejectsEmptyTeam(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS("", "frontend")
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Error("expected validation error for empty team")
	}
}

func TestValidator_RejectsInvalidTier(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS("test-team", "invalid-tier")
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Error("expected validation error for invalid tier")
	}
}

func TestValidator_AcceptsValidSpec(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS("test-team", "frontend")
	if err := v.ValidateCreate(context.Background(), ps); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
