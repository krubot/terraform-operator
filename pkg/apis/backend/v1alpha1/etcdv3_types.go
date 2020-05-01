package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="etcdv3s",singular="etcdv3",scope="Namespaced",shortName="etcdv3"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// Etcdv3 is the Schema for the Etcdv3s API
type EtcdV3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdV3Spec `json:"spec,omitempty"`
	Dep    []DepSpec  `json:"dep,omitempty"`
	Status StatusSpec `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// Etcdv3List contains a list of Etcdv3
type EtcdV3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EtcdV3 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdV3{}, &EtcdV3List{})
}
