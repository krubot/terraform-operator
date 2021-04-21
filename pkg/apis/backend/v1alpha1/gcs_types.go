package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:subresource:status
// GCSSpec defines the desired state of GCS
type GCSSpec struct {
	// GCS bucket name
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// GCS bucket prefix
	// +optional
	Prefix string `json:"prefix,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="gcs",singular="gcs",scope="Namespaced",shortName="gcs"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// GCS is the Schema for the GCS's API
type GCS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCSSpec    `json:"spec,omitempty"`
	Dep    []DepSpec  `json:"dep,omitempty"`
	Status StatusSpec `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GCSList contains a list of GCS
type GCSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []GCS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCS{}, &GCSList{})
}
