package stub

import (
  corev1 "k8s.io/api/core/v1"
)
type ParentType string

const (
	ParentPlan    = "tfplan"
	ParentApply   = "tfapply"
	ParentDestroy = "tfdestroy"
)

type TFPod struct {
	Image              string
	ImagePullPolicy    corev1.PullPolicy
	Namespace          string
	ProjectID          string
	Workspace          string
	Source             string
	ProviderConfigKeys map[string][]string
	BackendBucket      string
	BackendPrefix      string
	TFParent           string
	TFPlan             string
	TFInputs           map[string]string
	TFVars             map[string]string
	TFVarsFrom         map[string]string
}
