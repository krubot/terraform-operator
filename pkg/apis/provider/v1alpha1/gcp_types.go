package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:subresource:status
// GCPSpec defines the desired state of GCP
type GCPSpec struct {
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
// +kubebuilder:resource:path="gcp",singular="gcp",scope="Cluster",shortName="pro"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status",description="Description of the current status"
// GCP is the Schema for the GCPs API
type GCP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPSpec `json:"spec,omitempty"`
	Status string  `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// GCPList contains a list of GCP
type GCPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCP `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCP{}, &GCPList{})
}
