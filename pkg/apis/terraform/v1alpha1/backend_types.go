package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:subresource:status
// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// EtcdV3 backend endpoints
	// +optional
	Endpoints []string `json:"endpoints,omitempty"`
	// EtcdV3 backend lock
	// +optional
	Lock bool `json:"lock,omitempty"`
	// EtcdV3 backend prefix
	// +optional
	Prefix string `json:"prefix,omitempty"`
	// EtcdV3 backend cacert path
	// +optional
	CacertPath string `json:"cacert_path,omitempty"`
	// EtcdV3 backend cert path
	// +optional
	CertPath string `json:"cert_path,omitempty"`
	// EtcdV3 backend key path
	// +optional
	KeyPath string `json:"key_path,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="backends",singular="backend",scope="Cluster",shortName="bac"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status",description="Description of the current status"
// Backend is the Schema for the Backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec `json:"spec,omitempty"`
	Status string      `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
