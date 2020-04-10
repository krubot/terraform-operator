package controllers

import (
	"context"
	"reflect"
	"strings"

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
	module := &modulev1alpha1.GCS{}

	if err := r.Get(context.Background(), req.NamespacedName, module); err != nil {
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
		e := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
		h := map[string][]string{"Metadata-Flavor": {"Google"}}
		if ret := checkURL(e, h, 200); ret == nil {
			break
		}
	}

	if util.IsBeingDeleted(module) {
		finalizersProvider := &providerv1alpha1.Google{}
		finalizersModule := &modulev1alpha1.GCS{}
		finalizersBackend := &backendv1alpha1.EtcdV3{}

		for _, fin := range module.GetFinalizers() {
			instance_split_fin := strings.Split(fin, "_")

			if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GCS" && instance_split_fin[2] != module.ObjectMeta.Name {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersModule); errors.IsNotFound(err) {
					util.RemoveFinalizer(module, fin)

					if err := r.Update(context.Background(), module); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("module dependency is not met for deletion")
				}
			}

			if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersProvider); errors.IsNotFound(err) {
					util.RemoveFinalizer(module, fin)

					if err := r.Update(context.Background(), module); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("provider dependency is not met for deletion")
				}
			}

			if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersBackend); errors.IsNotFound(err) {
					util.RemoveFinalizer(module, fin)

					if err := r.Update(context.Background(), module); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("backend dependency is not met for deletion")
				}
			}
		}

		if err := terraform.WriteToFile([]byte("{}"), module.ObjectMeta.Namespace, "Module_"+module.Kind+"_"+module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformInit(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformNewWorkspace(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformSelectWorkspace(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformValidate(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformPlan(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.TerraformApply(module.ObjectMeta.Namespace); err != nil {
			return reconcile.Result{}, err
		}

		util.RemoveFinalizer(module, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)

		if err := r.Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// Set the initial depency state
	dependency_met := true

	// Over the list of dependencies
	for _, dep := range module.Dep {
		depProvider := &providerv1alpha1.Google{}
		depModule := &modulev1alpha1.GCS{}
		depBackend := &backendv1alpha1.EtcdV3{}

		if dep.Kind == "Backend" {
			if err := r.Get(context.Background(), types.NamespacedName{Name: dep.Name, Namespace: req.NamespacedName.Namespace}, depBackend); err != nil {
				return reconcile.Result{}, err
			}

			if depBackend.Status.State == "Success" {
				// Add finalizer to the module resource
				util.AddFinalizer(depBackend, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)

				// Update the CR with finalizer
				if err := r.Update(context.Background(), depBackend); err != nil {
					return reconcile.Result{}, err
				}
			} else {
				dependency_met = false
			}
		}

		if dep.Kind == "Module" {
			if err := r.Get(context.Background(), types.NamespacedName{Name: dep.Name, Namespace: req.NamespacedName.Namespace}, depModule); err != nil {
				return reconcile.Result{}, err
			}

			if depModule.Status.State == "Success" {
				// Add finalizer to the module resource
				util.AddFinalizer(depModule, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)

				// Update the CR with finalizer
				if err := r.Update(context.Background(), depModule); err != nil {
					return reconcile.Result{}, err
				}
			} else {
				dependency_met = false
			}
		}

		if dep.Kind == "Provider" {
			if err := r.Get(context.Background(), types.NamespacedName{Name: dep.Name, Namespace: req.NamespacedName.Namespace}, depProvider); err != nil {
				return reconcile.Result{}, err
			}

			if depProvider.Status.State == "Success" {
				// Add finalizer to the module resource
				util.AddFinalizer(depProvider, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)

				// Update the CR with finalizer
				if err := r.Update(context.Background(), depProvider); err != nil {
					return reconcile.Result{}, err
				}
			} else {
				dependency_met = false
			}
		}
	}

	// Check if dependency is met else interate again
	if !dependency_met {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Dependency"

		// Update the CR with status success
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Dependency", module.Status.Phase) {
		// Set the data
		module.Status.State = "Success"
		module.Status.Phase = "Dependency"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Add finalizer to the module resource
	util.AddFinalizer(module, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)

	// Update the CR with finalizer
	if err := r.Update(context.Background(), module); err != nil {
		return reconcile.Result{}, err
	}

	b, err := terraform.RenderModuleToTerraform(module.Spec, module.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = terraform.WriteToFile(b, module.ObjectMeta.Namespace, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual("Success", module.Status.State) {
		// Set the data
		module.Status.State = "Success"
		module.Status.Phase = "Output"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}
	}

	err = terraform.TerraformInit(module.ObjectMeta.Namespace)
	if err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Init"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if err := terraform.TerraformNewWorkspace(module.ObjectMeta.Namespace); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Workspace"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if err := terraform.TerraformSelectWorkspace(module.ObjectMeta.Namespace); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Workspace"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if err := terraform.TerraformValidate(module.ObjectMeta.Namespace); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Validate"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Validate", module.Status.Phase) {

		// Set the data
		module.Status.State = "Success"
		module.Status.Phase = "Validate"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err = terraform.TerraformPlan(module.ObjectMeta.Namespace); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Plan"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Plan", module.Status.Phase) {

		// Set the data
		module.Status.State = "Success"
		module.Status.Phase = "Plan"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err := terraform.TerraformApply(module.ObjectMeta.Namespace); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Apply"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Apply", module.Status.Phase) {

		// Set the data
		module.Status.State = "Success"
		module.Status.Phase = "Apply"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
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
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
