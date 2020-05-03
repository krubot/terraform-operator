package admission

import (
	"context"
	"net/http"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	modulev1alpha1 "github.com/krubot/terraform-operator/pkg/apis/module/v1alpha1"
	opa "github.com/krubot/terraform-operator/pkg/opa"
	terraform "github.com/krubot/terraform-operator/pkg/terraform"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
)

var path = "/tmp"

// terraformValidator validates GoogleStorageBucket
type TerraformValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// terraformValidator admits a pod iff a specific annotation exists.
func (v *TerraformValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	GoogleStorageBucket := &modulev1alpha1.GoogleStorageBucket{}
	GoogleStorageBucketIAMMember := &modulev1alpha1.GoogleStorageBucketIAMMember{}

	if req.Operation == admissionv1beta1.Delete {
		return admission.Allowed("")
	}

	if req.Kind.Kind == "GoogleStorageBucket" {
		err := v.decoder.Decode(req, GoogleStorageBucket)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		b, err := terraform.RenderModuleToTerraform(GoogleStorageBucket.Spec, strings.ToLower(GoogleStorageBucket.Kind)+"_"+strings.ToLower(GoogleStorageBucket.ObjectMeta.Name))
		if err != nil {
			return admission.Denied("Terraform render to module failed")
		}

		err = terraform.WriteToFile(b, GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)
		if err != nil {
			return admission.Denied("Terraform write to file failed")
		}

		if err := terraform.TerraformInit(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform init failed")
		}

		if err := terraform.TerraformNewWorkspace(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform new workspace failed")
		}

		if err := terraform.TerraformSelectWorkspace(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform select workspace failed")
		}

		if err := terraform.TerraformValidate(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform validate failed")
		}

		if err := terraform.TerraformPlan(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform plan failed")
		}

		if err := terraform.TerraformShow(GoogleStorageBucket.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

			return admission.Denied("Terraform plan failed")
		}

		terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name, path)

		if opa.Validation() != true {
			return admission.Denied("Terraform policies failed")
		}
	}

	if req.Kind.Kind == "GoogleStorageBucketIAMMember" {
		err := v.decoder.Decode(req, GoogleStorageBucketIAMMember)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		b, err := terraform.RenderModuleToTerraform(GoogleStorageBucketIAMMember.Spec, strings.ToLower(GoogleStorageBucketIAMMember.Kind)+"_"+strings.ToLower(GoogleStorageBucketIAMMember.ObjectMeta.Name))
		if err != nil {
			return admission.Denied("Terraform render to module failed")
		}

		err = terraform.WriteToFile(b, GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)
		if err != nil {
			return admission.Denied("Terraform write to file failed")
		}

		if err := terraform.TerraformInit(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform init failed")
		}

		if err := terraform.TerraformNewWorkspace(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform new workspace failed")
		}

		if err := terraform.TerraformSelectWorkspace(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform select workspace failed")
		}

		if err := terraform.TerraformValidate(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform validate failed")
		}

		if err := terraform.TerraformPlan(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform plan failed")
		}

		if err := terraform.TerraformShow(GoogleStorageBucketIAMMember.ObjectMeta.Namespace, path); err != nil {
			terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

			return admission.Denied("Terraform plan failed")
		}

		terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, path)

		if opa.Validation() != true {
			return admission.Denied("Terraform policies failed")
		}
	}

	return admission.Allowed("")
}

// terraformValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *TerraformValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
