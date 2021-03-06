package v1alpha1

// +kubebuilder:subresource:status
// StatusSpec status defines the status
type StatusSpec struct {
	// The current state of the terraform workflow
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	// +kubebuilder:validation:Enum={"Dependency","Output","Init","Workspace","Validate","Plan","Apply"}
	Phase string `json:"phase"`
}
