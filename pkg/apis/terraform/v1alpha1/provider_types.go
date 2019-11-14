package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:subresource:status
// ProviderSpec defines the desired state of Provider
type ProviderSpec struct {
	// Kubernetes provider hostname (in form of URI) of Kubernetes master
	// +optional
	Host string `json:"host,omitempty"`
	// Kubernetes provider token of your service account
	// +optional
	Token string `json:"token,omitempty"`
	// Kubernetes provider whether server should be accessed without verifying the TLS certificate
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=providers,singular=provider,scope=Cluster,shortName=provider
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status",description="Description of the current status"
// Provider is the Schema for the providers API
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec `json:"spec,omitempty"`
	Status string       `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Provider{}, &ProviderList{})
}
