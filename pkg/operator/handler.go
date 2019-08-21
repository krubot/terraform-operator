package stub

import (
	"fmt"
	"context"

	"github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
	"github.com/krubot/terraform-operator/pkg/terraform"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

const ResourceName = "aws_s3_bucket"

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct{}

// ObjectCreated is called when an object is created
func (t *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch tf := event.Object.(type) {
	case *v1alpha1.AwsS3Bucket:
		tf.Status = "Started"
		err := sdk.Update(tf)
		if err != nil {
			return err
		}
		b, err := terraform.RenderToTerraform(tf.Spec, ResourceName, tf.metadata.name)
		if err != nil {
			return err
		}
		logrus.Infof("%s", string(b))
		err = terraform.WriteToFile(b, fmt.Sprintf("%s-%s", ResourceName, uid))
		if err != nil {
			return err
		}
		tf.Status = "Created"
		err = sdk.Update(tf)
		if err != nil {
			return err
		}
		err = terraform.TerraformValidate()
		if err != nil {
			return err
		}
		tf.Status = "Validated"
		err = sdk.Update(tf)
		if err != nil {
			return err
		}
	}
	return nil
}
