package controllers

import (
	"context"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	backendv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/backend/v1alpha1"
	modulev1alpha1 "github.com/krubot/terraform-operator/pkg/apis/module/v1alpha1"
	providerv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/provider/v1alpha1"
	terraform "github.com/krubot/terraform-operator/pkg/terraform"
	util "github.com/krubot/terraform-operator/pkg/util"
)

// ReconcileProvider reconciles a Backend object
type ReconcileModule struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

func (r *ReconcileModule) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// Fetch the Module instance
	instance := &modulev1alpha1.GCS{}

	err := r.Get(context.Background(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	for {
		backend := &backendv1alpha1.EtcdV3{}
		r.Get(context.Background(), types.NamespacedName{Name: "state", Namespace: ""}, backend)

		provider := &providerv1alpha1.GCP{}
		r.Get(context.Background(), types.NamespacedName{Name: "cloud", Namespace: ""}, provider)

		if backend.Status == "Ready" && provider.Status == "Ready" {
			break
		}

		// Wait before loop
		time.Sleep(1 * time.Second)
	}

	b, err := terraform.RenderModuleToTerraform(instance.Spec, instance.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = terraform.WriteToFile(b, instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual("Success", instance.Status.State) {
		// Add finalizer to the module resource
		util.AddFinalizer(instance, "controller_module")

		// Update the CR with finalizer
		if err := r.Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		// Set the data
		instance.Status.State = "Success"
		instance.Status.Phase = "Output"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	if util.IsBeingDeleted(instance) {
		if !util.HasFinalizer(instance, "controller_module") {
			return reconcile.Result{}, nil
		}

		err = terraform.WriteToFile([]byte("{}"), instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformInit(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformNewWorkspace(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformSelectWorkspace(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformValidate(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformPlan(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.TerraformApply(instance.ObjectMeta.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		util.RemoveFinalizer(instance, "controller_module")

		err = r.Update(context.Background(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	err = terraform.TerraformInit(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Init"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// Wait before loop
	time.Sleep(1 * time.Second)

	err = terraform.TerraformNewWorkspace(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Workspace"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// Wait before loop
	time.Sleep(1 * time.Second)

	err = terraform.TerraformSelectWorkspace(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Workspace"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// Wait before loop
	time.Sleep(1 * time.Second)

	err = terraform.TerraformValidate(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Validate"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Validate", instance.Status.Phase) {

		// Set the data
		instance.Status.State = "Success"
		instance.Status.Phase = "Validate"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Wait before loop
	time.Sleep(1 * time.Second)

	err = terraform.TerraformPlan(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Plan"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Plan", instance.Status.Phase) {

		// Set the data
		instance.Status.State = "Success"
		instance.Status.Phase = "Plan"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Wait before loop
	time.Sleep(1 * time.Second)

	err = terraform.TerraformApply(instance.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		instance.Status.State = "Failure"
		instance.Status.Phase = "Apply"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(instance.ObjectMeta.Namespace, instance.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Apply", instance.Status.Phase) {

		// Set the data
		instance.Status.State = "Success"
		instance.Status.Phase = "Apply"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileModule) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.GCS{}).
		Watches(&source.Kind{Type: &modulev1alpha1.GCS{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationOrFinalizerChangedPredicate{}).
		Complete(r)
}
