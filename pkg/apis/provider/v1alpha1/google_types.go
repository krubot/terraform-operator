package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:subresource:status
// GoogleSpec defines the desired state of Google
type GoogleSpec struct {
	// Either the path to or the contents of a service account key file in JSON format.
	// +optional
	Credentials string `json:"credentials,omitempty"`
	// The default project to manage resources in. If another project is specified on a resource, it will take precedence.
	// +optional
	Project string `json:"project,omitempty"`
	// The default region to manage resources in. If another region is specified on a regional resource, it will take precedence.
	// +optional
	Region string `json:"region,omitempty"`
	// The default zone to manage resources in. Generally, this zone should be within the default region you specified. If another zone is specified on a zonal resource, it will take precedence.
	// +optional
	Zone string `json:"zone,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path="googles",singular="google",scope="Namespaced",shortName="google"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// Google is the Schema for the Googles API
type Google struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleSpec `json:"spec,omitempty"`
	Dep    []DepSpec  `json:"dep,omitempty"`
	Status StatusSpec `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GoogleList contains a list of Google
type GoogleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Google `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Google{}, &GoogleList{})
}
