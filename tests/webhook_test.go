package webhook_test

import (
	"context"
	"testing"

	"github.com/SumitDalavi/k8s-golden-path-provisioner/webhook"
	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newPS(replicas int, port int32, image string) *platformv1.PlatformService {
	return &platformv1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec:       platformv1.PlatformServiceSpec{Replicas: replicas, Port: port, Image: image},
	}
}

func TestDefaulter_SetsReplicaDefault(t *testing.T) {
	d := &webhook.PlatformServiceDefaulter{}
	ps := newPS(0, 0, "")
	if err := d.Default(context.Background(), ps); err != nil {
		t.Fatal(err)
	}
	if ps.Spec.Replicas != 2 { t.Errorf("expected replicas=2, got %d", ps.Spec.Replicas) }
	if ps.Spec.Port != 8080 { t.Errorf("expected port=8080, got %d", ps.Spec.Port) }
	if ps.Labels["app.kubernetes.io/managed-by"] != "golden-path-provisioner" {
		t.Error("missing managed-by label")
	}
}

func TestValidator_RejectsZeroReplicas(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS(0, 8080, "nginx:1.25")
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Error("expected validation error for 0 replicas")
	}
}

func TestValidator_RejectsTooManyReplicas(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS(100, 8080, "nginx:1.25")
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Error("expected validation error for 100 replicas")
	}
}

func TestValidator_AcceptsValidSpec(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS(3, 9090, "myapp:v1")
	if err := v.ValidateCreate(context.Background(), ps); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidator_RejectsEmptyImage(t *testing.T) {
	v := &webhook.PlatformServiceValidator{}
	ps := newPS(2, 8080, "")
	if err := v.ValidateCreate(context.Background(), ps); err == nil {
		t.Error("expected validation error for empty image")
	}
}
