package stub

import (
	"context"

	"github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
	"github.com/krubot/terraform-operator/pkg/terraform"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

const ResourceName="aws_s3_bucket"

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {}

// ObjectCreated is called when an object is created
func (t *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.AwsS3Bucket:
		b, err := terraform.RenderToTerraform(o.Spec, ResourceName, string(o.GetUID()))
		if err != nil {
			return err
		}
		logrus.Infof("%s", string(b))
	}
	return nil
}
