package v1alpha1

// +kubebuilder:subresource:status
// Google status defines the status of Google
type StatusSpec struct {
	// The current state of the terraform workflow
	// +kubebuilder:validation:Enum={"Success","Failure"}
	State string `json:"state"`
	// The current phase of the terraform workflow
	// +kubebuilder:validation:Enum={"Dependency","Output","Init","Workspace","Validate","Plan","Apply"}
	Phase string `json:"phase"`
}
