package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

// +kubebuilder:subresource:status
// Google status defines the status of Google
type GoogleStatus struct {
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	Phase string `json:"phase"`
}

// +kubebuilder:subresource:status
// DepSpec defines the dependency list of Google
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
// +kubebuilder:resource:path="google",singular="google",scope="Namespaced",shortName="pro"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Description of the current state"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Description of the current phase"
// Google is the Schema for the Googles API
type Google struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleSpec   `json:"spec,omitempty"`
	Dep    []DepSpec    `json:"dep,omitempty"`
	Status GoogleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GoogleList contains a list of Google
type GoogleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Google `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Google{}, &GoogleList{})
}
