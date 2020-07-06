package v1alpha1

// +kubebuilder:subresource:status
// DepSpec defines the dependency list of Google
type DepSpec struct {
	// Dependency kind
	// +kubebuilder:validation:Enum={"Backend","Module","Provider"}
	Kind string `json:"kind"`
	// Dependency name
	Name string `json:"name"`
	// Dependency type
	// +kubebuilder:validation:Enum={"EtcdV3","GCS","GoogleStorageBucket","GoogleStorageBucketIAMMember","Google"}
	Type string `json:"type"`
}
