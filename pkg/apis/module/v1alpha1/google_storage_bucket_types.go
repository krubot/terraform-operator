package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:subresource:spec
// GoogleStorageBucketSpec defines the desired state of GoogleStorageBucket
type GoogleStorageBucketSpec struct {
	// Kubernetes namespace GoogleStorageBucket source
	// +kubebuilder:validation:Enum={"/opt/modules/gcp/google_storage_bucket/"}
	Source string `json:"source"`
	// GCP bucket name
	Name string `json:"name"`
}

// +kubebuilder:subresource:spec
// OutputSpec defines the desired state of Output
type OutputSpec struct {
	// Bucket name (for single use)
	Name string `json:"name"`
	// Bucket URL (for single use)
	URL string `json:"url"`
}

// +genclient
// +genclient:Namespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="googlestoragebuckets",singular="googlestoragebucket",scope="Namespaced",shortName="gsb"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// GoogleStorageBucket is the Schema for the GoogleStorageBuckets API
type GoogleStorageBucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleStorageBucketSpec `json:"spec,omitempty"`
	Dep    []DepSpec               `json:"dep,omitempty"`
	Output OutputSpec              `json:"output,omitempty"`
	Status StatusSpec              `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GoogleStorageBucketList contains a list of GoogleStorageBucket
type GoogleStorageBucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []GoogleStorageBucket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleStorageBucket{}, &GoogleStorageBucketList{})
}
