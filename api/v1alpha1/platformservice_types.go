package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PlatformServiceSpec defines the desired state of a PlatformService
type PlatformServiceSpec struct {
	// Team is the owning team of this service (used for tagging and RBAC)
	Team string `json:"team"`

	// Tier defines the architecture tier (frontend, backend, worker)
	// +kubebuilder:validation:Enum=frontend;backend;worker
	Tier string `json:"tier"`

	// ExposeExternally determines if an external Ingress/Gateway should be provisioned
	ExposeExternally bool `json:"exposeExternally,omitempty"`
}

// PlatformServiceStatus defines the observed state of PlatformService
type PlatformServiceStatus struct {
	// Phase represents the current state of the environment provisioning
	Phase string `json:"phase,omitempty"`

	// Namespace is the name of the provisioned namespace
	Namespace string `json:"namespace,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Team",type=string,JSONPath=`.spec.team`
// +kubebuilder:printcolumn:name="Tier",type=string,JSONPath=`.spec.tier`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`

// PlatformService is the Schema for the platformservices API
type PlatformService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformServiceSpec   `json:"spec,omitempty"`
	Status PlatformServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type PlatformServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlatformService `json:"items"`
}
