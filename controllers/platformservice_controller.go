package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PlatformServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PlatformServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var ps unstructured.Unstructured
	ps.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "platform.example.com",
		Version: "v1alpha1",
		Kind:    "PlatformService",
	})
	
	if err := r.Get(ctx, req.NamespacedName, &ps); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	nsName := fmt.Sprintf("%s-prod", ps.GetName())

	// 1. Reconcile Namespace
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName}}
	if err := r.Create(ctx, ns); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}
	logger.Info("Ensured Namespace", "namespace", nsName)

	// 2. Reconcile RBAC RoleBinding
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cicd-deployer-binding",
			Namespace: nsName,
		},
		Subjects: []rbacv1.Subject{{Kind: "ServiceAccount", Name: "default", Namespace: nsName}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "edit",
		},
	}
	if err := r.Create(ctx, rb); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// 3. Reconcile ResourceQuota
	rq := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-quota",
			Namespace: nsName,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			},
		},
	}
	if err := r.Create(ctx, rq); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PlatformServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "platform.example.com",
		Version: "v1alpha1",
		Kind:    "PlatformService",
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(u).
		Complete(r)
}
