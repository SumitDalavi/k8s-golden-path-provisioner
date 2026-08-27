package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1alpha1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
)

type PlatformServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PlatformServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var ps platformv1alpha1.PlatformService
	if err := r.Get(ctx, req.NamespacedName, &ps); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	nsName := fmt.Sprintf("svc-%s", ps.Name)

	// 1. Reconcile Namespace
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName}}
	if err := r.Create(ctx, ns); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}
	logger.Info("Ensured Namespace", "namespace", nsName)

	// 2. Reconcile NetworkPolicy (Default Deny Ingress)
	netpol := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny-ingress",
			Namespace: nsName,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
		},
	}
	if err := r.Create(ctx, netpol); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// 3. Reconcile CI/CD ServiceAccount
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cicd-deployer",
			Namespace: nsName,
		},
	}
	if err := r.Create(ctx, sa); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// 4. Reconcile RoleBinding for CI/CD
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cicd-deployer-binding",
			Namespace: nsName,
		},
		Subjects: []rbacv1.Subject{{Kind: "ServiceAccount", Name: sa.Name, Namespace: nsName}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "edit", // Standard K8s role
		},
	}
	if err := r.Create(ctx, rb); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// 5. Update Status
	ps.Status.Phase = "Ready"
	ps.Status.Namespace = nsName
	if err := r.Status().Update(ctx, &ps); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PlatformServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.PlatformService{}).
		Complete(r)
}
