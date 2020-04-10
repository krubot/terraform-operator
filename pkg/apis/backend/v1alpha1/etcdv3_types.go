package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:subresource:status
// Etcdv3Spec defines the desired state of Etcdv3
type EtcdV3Spec struct {
	// EtcdV3 Etcdv3 endpoints
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

// +kubebuilder:subresource:status
// EtcdV3 status defines the status of EtcdV3
type EtcdV3Status struct {
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	Phase string `json:"phase"`
}

// +kubebuilder:subresource:status
// DepSpec defines the dependency list of EtcdV3
type DepSpec struct {
	// Dependency kind
	// +kubebuilder:validation:Enum={"Backend","Module","Provider"}
	Kind string `json:"kind"`
	// Dependency name
	Name string `json:"name"`
	// Dependency type
	// +kubebuilder:validation:Enum={"EtcdV3","GCS","Google"}
	Type string `json:"type"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="etcdv3",singular="etcdv3",scope="Namespaced",shortName="bac"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// Etcdv3 is the Schema for the Etcdv3s API
type EtcdV3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdV3Spec   `json:"spec,omitempty"`
	Dep    []DepSpec    `json:"dep,omitempty"`
	Status EtcdV3Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// Etcdv3List contains a list of Etcdv3
type EtcdV3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdV3 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdV3{}, &EtcdV3List{})
}
