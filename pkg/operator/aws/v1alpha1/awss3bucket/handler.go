package awss3bucket

import (
    log "github.com/Sirupsen/logrus"
    "github.com/krubot/terraform-operator/pkg/terraform"
    "github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
)

const ResourceName="aws_s3_bucket"
type Handler struct{}

// Init is used for initialization logic
func (t *Handler) Init() error {
	return nil
}

// ObjectCreated is called when an object is created
func (t *Handler) ObjectCreated(obj interface{}) {
  o := obj.(*v1alpha1.AwsS3Bucket)
	b, err := terraform.RenderToTerraform(o.Spec, ResourceName, string(o.GetUID()))
  
	if err != nil {
		log.Info(err)
	}

	log.Infof("%s", string(b))
}

// ObjectDeleted is called when an object is deleted
func (t *Handler) ObjectDeleted(obj interface{}) {
}

// ObjectUpdated is called when an object is updated
func (t *Handler) ObjectUpdated(objOld, objNew interface{}) {
}
