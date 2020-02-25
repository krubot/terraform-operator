package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:subresource:spec
// GCSSpec defines the desired state of GCS
type GCSSpec struct {
	// Kubernetes namespace GCS source
	// +kubebuilder:validation:Enum={"/var/lib/modules/gcp/gcs/"}
	Source string `json:"source"`
	// GCP bucket name
	Name string `json:"name"`
	// GCP bucket project
	Project string `json:"project"`
}

// +kubebuilder:subresource:status
// GCS status defines the status of GCS
type GCSStatus struct {
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	Phase string `json:"phase"`
}

// +genclient
// +genclient:Namespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="gcs",singular="gcs",scope="Namespaced",shortName="mod"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// GCS is the Schema for the GCSs API
type GCS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCSSpec   `json:"spec,omitempty"`
	Status GCSStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GCSList contains a list of GCS
type GCSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCS{}, &GCSList{})
}
