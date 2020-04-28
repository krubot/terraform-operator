package controllers

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"

	modulev1alpha1 "github.com/krubot/terraform-operator/pkg/apis/module/v1alpha1"
	terraform "github.com/krubot/terraform-operator/pkg/terraform"
)

func (r *ReconcileModule) terraformReconcile(moduleInterface interface{}) error {
	switch module := moduleInterface.(type) {
	case *modulev1alpha1.GCS:
		if err := terraform.TerraformInit(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Init"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if !reflect.DeepEqual("Init", module.Status.Phase) {
			// Set the data
			module.Status.State = "Success"
			module.Status.Phase = "Init"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
		}

		if err := terraform.TerraformNewWorkspace(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Workspace"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if err := terraform.TerraformSelectWorkspace(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Workspace"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if !reflect.DeepEqual("Workspace", module.Status.Phase) {
			// Set the data
			module.Status.State = "Success"
			module.Status.Phase = "Workspace"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
		}

		if err := terraform.TerraformValidate(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Validate"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if !reflect.DeepEqual("Validate", module.Status.Phase) {
			// Set the data
			module.Status.State = "Success"
			module.Status.Phase = "Validate"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
		}

		if err := terraform.TerraformPlan(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Plan"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if !reflect.DeepEqual("Plan", module.Status.Phase) {
			// Set the data
			module.Status.State = "Success"
			module.Status.Phase = "Plan"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
		}

		if err := terraform.TerraformApply(module.ObjectMeta.Namespace); err != nil {
			// Set the data
			module.Status.State = "Failure"
			module.Status.Phase = "Apply"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
			if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
				return err
			}
			return nil
		}

		if !reflect.DeepEqual("Apply", module.Status.Phase) {
			// Set the data
			module.Status.State = "Success"
			module.Status.Phase = "Apply"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), module); err != nil {
				return err
			}
		}
	default:
		return errors.NewBadRequest("terraform reconcile does not implement on this type")
	}

	return nil
}
