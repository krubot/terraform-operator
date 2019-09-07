package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProviderSpec defines the desired state of Provider
// +kubebuilder:subresource:status
type ProviderSpec struct {
		// Kubernetes provider hostname (in form of URI) of Kubernetes master
		Host			string	`json:"host,omitempty"`
		// Kubernetes provider token of your service account
		Token			string	`json:"token,omitempty"`
		// Kubernetes provider whether server should be accessed without verifying the TLS certificate
		Insecure	bool		`json:"insecure,omitempty"`

		// AWS provider region
		// Region			string					`json:"region,omitempty"`
		// AWS provider assume role
		// AssumeRole	AssumeRoleSpec	`json:"assume_role,omitempty"`
}

// type AssumeRoleSpec struct {
// 		RoleArn			string	`json:"role_arn,omitempty"`
// 		SessionName	string	`json:"session_name,omitempty"`
// 		ExternalId	string	`json:"external_id,omitempty"`
// 		Policy			string	`json:"policy,omitempty"`
// }

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Provider is the Schema for the providers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec `json:"spec,omitempty"`
	Status string     	`json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Provider{}, &ProviderList{})
}
