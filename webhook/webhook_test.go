package webhook

import (
	"context"
	"testing"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaulter_SetsDefaults(t *testing.T) {
	d := &PlatformServiceDefaulter{}
	ps := &platformv1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}
	if err := d.Default(context.Background(), ps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.Spec.Tier != "backend" {
		t.Errorf("expected default tier 'backend', got %s", ps.Spec.Tier)
	}
	if ps.Spec.Team != "default-team" {
		t.Errorf("expected default team 'default-team', got %s", ps.Spec.Team)
	}
	if ps.Labels["app.kubernetes.io/managed-by"] != "golden-path-provisioner" {
		t.Errorf("expected managed-by label to be set")
	}
}

func TestDefaulter_DoesNotOverride(t *testing.T) {
	d := &PlatformServiceDefaulter{}
	ps := &platformv1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: platformv1.PlatformServiceSpec{
			Tier: "frontend",
			Team: "my-team",
		},
	}
	if err := d.Default(context.Background(), ps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.Spec.Tier != "frontend" {
		t.Errorf("expected tier 'frontend' to be preserved, got %s", ps.Spec.Tier)
	}
	if ps.Spec.Team != "my-team" {
		t.Errorf("expected team 'my-team' to be preserved, got %s", ps.Spec.Team)
	}
}

func TestDefaulter_NonPlatformService(t *testing.T) {
	d := &PlatformServiceDefaulter{}
	err := d.Default(context.Background(), &corev1.ConfigMap{})
	if err == nil {
		t.Fatalf("expected error for non-PlatformService object")
	}
}

func TestDefaulter_SetsLabels(t *testing.T) {
	d := &PlatformServiceDefaulter{}
	ps := &platformv1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test",
			Labels: map[string]string{"existing": "label"},
		},
		Spec: platformv1.PlatformServiceSpec{Tier: "worker", Team: "ops"},
	}
	if err := d.Default(context.Background(), ps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.Labels["existing"] != "label" {
		t.Errorf("expected existing label to be preserved")
	}
	if ps.Labels["app.kubernetes.io/managed-by"] != "golden-path-provisioner" {
		t.Errorf("expected managed-by label to be set")
	}
}

func TestValidator_ValidCreate(t *testing.T) {
	v := &PlatformServiceValidator{}
	ps := &platformv1.PlatformService{
		Spec: platformv1.PlatformServiceSpec{Tier: "backend", Team: "my-team"},
	}
	if err := v.ValidateCreate(context.Background(), ps); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidator_InvalidTier(t *testing.T) {
	v := &PlatformServiceValidator{}
	ps := &platformv1.PlatformService{
		Spec: platformv1.PlatformServiceSpec{Tier: "invalid", Team: "my-team"},
	}
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Errorf("expected error for invalid tier")
	}
}

func TestValidator_MissingTeam(t *testing.T) {
	v := &PlatformServiceValidator{}
	ps := &platformv1.PlatformService{
		Spec: platformv1.PlatformServiceSpec{Tier: "backend"},
	}
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Errorf("expected error for missing team")
	}
}

func TestValidator_ValidUpdate(t *testing.T) {
	v := &PlatformServiceValidator{}
	old := &platformv1.PlatformService{
		Spec: platformv1.PlatformServiceSpec{Tier: "backend", Team: "team-a"},
	}
	new := &platformv1.PlatformService{
		Spec: platformv1.PlatformServiceSpec{Tier: "worker", Team: "team-b"},
	}
	if err := v.ValidateUpdate(context.Background(), old, new); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidator_ValidDelete(t *testing.T) {
	v := &PlatformServiceValidator{}
	ps := &platformv1.PlatformService{}
	if err := v.ValidateDelete(context.Background(), ps); err != nil {
		t.Errorf("unexpected error on delete: %v", err)
	}
}

func TestValidator_NonPlatformService(t *testing.T) {
	v := &PlatformServiceValidator{}
	err := v.ValidateCreate(context.Background(), &corev1.ConfigMap{})
	if err == nil {
		t.Errorf("expected error for non-PlatformService object")
	}
}

func TestValidator_AllValidTiers(t *testing.T) {
	v := &PlatformServiceValidator{}
	for _, tier := range []string{"frontend", "backend", "worker"} {
		ps := &platformv1.PlatformService{
			Spec: platformv1.PlatformServiceSpec{Tier: tier, Team: "team"},
		}
		if err := v.ValidateCreate(context.Background(), ps); err != nil {
			t.Errorf("unexpected error for tier %s: %v", tier, err)
		}
	}
}
