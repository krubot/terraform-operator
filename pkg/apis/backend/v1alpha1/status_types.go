package v1alpha1

// +kubebuilder:subresource:status
// EtcdV3 status defines the status of EtcdV3
type StatusSpec struct {
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	Phase string `json:"phase"`
}
