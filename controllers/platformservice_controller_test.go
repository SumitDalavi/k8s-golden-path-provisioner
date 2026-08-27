package controllers

import (
	"context"
	"testing"

	platformv1alpha1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = platformv1alpha1.AddToScheme(s)
	_ = corev1.AddToScheme(s)
	_ = networkingv1.AddToScheme(s)
	_ = rbacv1.AddToScheme(s)
	return s
}

func TestReconcile_NotFound(t *testing.T) {
	s := newScheme()
	fakeClient := fake.NewClientBuilder().WithScheme(s).Build()

	r := &PlatformServiceReconciler{Client: fakeClient, Scheme: s}
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "nonexistent", Namespace: "default"}}

	res, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Requeue {
		t.Errorf("did not expect requeue")
	}
}

func TestReconcile_Success(t *testing.T) {
	s := newScheme()

	ps := &platformv1alpha1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-service",
			Namespace: "default",
		},
		Spec: platformv1alpha1.PlatformServiceSpec{
			Team: "team-a",
			Tier: "backend",
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(s).
		WithStatusSubresource(&platformv1alpha1.PlatformService{}).
		WithObjects(ps).
		Build()

	r := &PlatformServiceReconciler{Client: fakeClient, Scheme: s}
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "my-service", Namespace: "default"}}

	_, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the namespace was created
	ns := &corev1.Namespace{}
	if err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "svc-my-service"}, ns); err != nil {
		t.Errorf("namespace svc-my-service was not created: %v", err)
	}

	// Verify the NetworkPolicy was created
	netpol := &networkingv1.NetworkPolicy{}
	if err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "default-deny-ingress", Namespace: "svc-my-service"}, netpol); err != nil {
		t.Errorf("NetworkPolicy was not created: %v", err)
	}

	// Verify the ServiceAccount was created
	sa := &corev1.ServiceAccount{}
	if err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "cicd-deployer", Namespace: "svc-my-service"}, sa); err != nil {
		t.Errorf("ServiceAccount was not created: %v", err)
	}

	// Verify the RoleBinding was created
	rb := &rbacv1.RoleBinding{}
	if err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "cicd-deployer-binding", Namespace: "svc-my-service"}, rb); err != nil {
		t.Errorf("RoleBinding was not created: %v", err)
	}
}

func TestReconcile_AlreadyExists(t *testing.T) {
	s := newScheme()

	ps := &platformv1alpha1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "my-service", Namespace: "default"},
		Spec:       platformv1alpha1.PlatformServiceSpec{Team: "team-a", Tier: "backend"},
	}
	// Pre-create resources that reconciler will try to create
	existingNs := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "svc-my-service"}}

	fakeClient := fake.NewClientBuilder().
		WithScheme(s).
		WithStatusSubresource(&platformv1alpha1.PlatformService{}).
		WithObjects(ps, existingNs).
		Build()

	r := &PlatformServiceReconciler{Client: fakeClient, Scheme: s}
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "my-service", Namespace: "default"}}

	_, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error on already-exists: %v", err)
	}
}

func TestReconcileMultipleTimes(t *testing.T) {
	s := newScheme()
	ps := &platformv1alpha1.PlatformService{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-multi", Namespace: "default"},
		Spec:       platformv1alpha1.PlatformServiceSpec{Team: "team", Tier: "frontend"},
	}
	fakeClient := fake.NewClientBuilder().
		WithScheme(s).
		WithStatusSubresource(&platformv1alpha1.PlatformService{}).
		WithObjects(ps).
		Build()
	r := &PlatformServiceReconciler{Client: fakeClient, Scheme: s}
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "svc-multi", Namespace: "default"}}

	// First reconcile
	_, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("first reconcile error: %v", err)
	}
	// Second reconcile (all resources already exist - AlreadyExists path)
	_, err = r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("second reconcile error: %v", err)
	}
}
