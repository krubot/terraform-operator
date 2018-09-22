package stub

import (
	"context"

	"github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/krubot/terraform-operator/pkg/terraform"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{
		controller: terraform.NewController()
	}
}

type Handler struct {
	controller	*[]terraform.Controller
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Terraform:
		sdk.Create(newTerraformPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("failed to create busybox pod : %v", err)
			return err
		}
	}
	return nil
}
