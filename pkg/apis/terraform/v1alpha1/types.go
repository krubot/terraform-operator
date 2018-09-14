package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Terraform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TerraformSpec           `json:"spec,omitempty"`
	Status            TerraformOperatorStatus `json:"status"`
}

// TerraformSpec is the top level structure of the spec body
type TerraformSpec struct {
	Image           string                                 `json:"image",omitempty`
	BackendBucket   string                                 `json:"backendBucket,omitempty"`
	ProviderConfig  map[string]TerraformSpecProviderConfig `json:"providerConfig,omitempty"`
	TFSource        string                                 `json:"sources,omitempty"`
	TFInputs        []TerraformConfigInputs                `json:"tfinputs,omitempty"`
	TFVars          map[string]string                      `json:"tfvars,omitempty"`
}

type TerraformSpecProviderConfig struct {
	SecretName string `json:"secretName,omitempty"`
}

type TerraformConfigInputs struct {
	Name   string            `json:"name,omitempty"`
	VarMap map[string]string `json:"varMap,omitempty"`
}

// To be defined later
type Status struct {}
