package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:subresource:spec
// GoogleStorageBucketIAMMemberSpec defines the desired state of GoogleStorageBucketIAMMember
type GoogleStorageBucketIAMMemberSpec struct {
	// Kubernetes namespace GoogleStorageBucketIAMMember source
	// +kubebuilder:validation:Enum={"/opt/modules/gcp/google_storage_bucket_iam_member/"}
	Source string `json:"source"`
	// Map of role (key) and list of members (value) to add the IAM policies/bindings
	Bindings map[string][]string `json:"bindings"`
	// Entities list to add the IAM policies/bindings
	Entities []string `json:"entities"`
}

// +kubebuilder:subresource:spec
// OutputGoogleStorageBucketSpec defines the desired state of Output
type OutputGoogleStorageBucketIAMMemberSpec struct{}

// +genclient
// +genclient:Namespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="googlestoragebucketiammembers",singular="googlestoragebucketiammember",scope="Namespaced",shortName="gsbim"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// GoogleStorageBucketIAMMember is the Schema for the GoogleStorageBucketIAMMembers API
type GoogleStorageBucketIAMMember struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleStorageBucketIAMMemberSpec       `json:"spec,omitempty"`
	Dep    []DepSpec                              `json:"dep,omitempty"`
	Output OutputGoogleStorageBucketIAMMemberSpec `json:"output,omitempty"`
	Status StatusSpec                             `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GoogleStorageBucketIAMMemberList contains a list of GoogleStorageBucketIAMMember
type GoogleStorageBucketIAMMemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []GoogleStorageBucketIAMMember `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleStorageBucketIAMMember{}, &GoogleStorageBucketIAMMemberList{})
}
