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

func (r *ReconcileModule) deletionReconcile(moduleInterface interface{}, finalizerInterfaces ...interface{}) error {
	for _, finalizerInterface := range finalizerInterfaces {
		switch module := moduleInterface.(type) {
		case *modulev1alpha1.GoogleStorageBucket:
			for _, fin := range module.GetFinalizers() {
				instance_split_fin := strings.Split(fin, "_")
				switch finalizer := finalizerInterface.(type) {
				case *modulev1alpha1.GoogleStorageBucket:
					if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucket" && instance_split_fin[2] != module.ObjectMeta.Name {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("GoogleStorageBucket dependency is not met for deletion")
						}
					}
				case *modulev1alpha1.GoogleStorageBucketIAMMember:
					if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucketIAMMember" && instance_split_fin[2] != module.ObjectMeta.Name {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("GoogleStorageBucketIAMMember dependency is not met for deletion")
						}
					}
				case *providerv1alpha1.Google:
					if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("Google dependency is not met for deletion")
						}
					}
				case *backendv1alpha1.EtcdV3:
					if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("EtcdV3 dependency is not met for deletion")
						}
					}
				}
			}
		case *modulev1alpha1.GoogleStorageBucketIAMMember:
			for _, fin := range module.GetFinalizers() {
				instance_split_fin := strings.Split(fin, "_")
				switch finalizer := finalizerInterface.(type) {
				case *modulev1alpha1.GoogleStorageBucket:
					if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucket" && instance_split_fin[2] != module.ObjectMeta.Name {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("GoogleStorageBucket dependency is not met for deletion")
						}
					}
				case *modulev1alpha1.GoogleStorageBucketIAMMember:
					if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucketIAMMember" && instance_split_fin[2] != module.ObjectMeta.Name {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("GoogleStorageBucketIAMMember dependency is not met for deletion")
						}
					}
				case *providerv1alpha1.Google:
					if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("Google dependency is not met for deletion")
						}
					}
				case *backendv1alpha1.EtcdV3:
					if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(module, fin)
							if err := r.Update(context.Background(), module); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("EtcdV3 dependency is not met for deletion")
						}
					}
				}
			}
		}
	}
	return nil
}

func (r *ReconcileModule) dependencyReconcile(moduleInterface interface{}, depInterfaces ...interface{}) (bool, error) {
	// Set the initial depency state
	dependency_met := true
	// Over the list of dependencies
	for _, depInterface := range depInterfaces {
		switch module := moduleInterface.(type) {
		case *modulev1alpha1.GoogleStorageBucket:
			for _, depModule := range module.Dep {
				switch dep := depInterface.(type) {
				case *modulev1alpha1.GoogleStorageBucket:
					if depModule.Kind == "Module" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *modulev1alpha1.GoogleStorageBucketIAMMember:
					if depModule.Kind == "Module" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *providerv1alpha1.Google:
					if depModule.Kind == "Provider" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *backendv1alpha1.EtcdV3:
					if depModule.Kind == "Backend" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				}
			}
		case *modulev1alpha1.GoogleStorageBucketIAMMember:
			for _, depModule := range module.Dep {
				switch dep := depInterface.(type) {
				case *modulev1alpha1.GoogleStorageBucket:
					if depModule.Kind == "Module" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *modulev1alpha1.GoogleStorageBucketIAMMember:
					if depModule.Kind == "Module" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *providerv1alpha1.Google:
					if depModule.Kind == "Provider" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *backendv1alpha1.EtcdV3:
					if depModule.Kind == "Backend" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							util.AddFinalizer(dep, "Module_"+module.Kind+"_"+module.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				}
			}
		}
	}
	return dependency_met, nil
}

func (r *ReconcileModule) terraformReconcile(moduleInterface interface{}) error {
	switch module := moduleInterface.(type) {
	case *modulev1alpha1.GoogleStorageBucket:
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
	case *modulev1alpha1.GoogleStorageBucketIAMMember:
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

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

func (r *ReconcileModule) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	GoogleStorageBucket := &modulev1alpha1.GoogleStorageBucket{}
	GoogleStorageBucketIAMMember := &modulev1alpha1.GoogleStorageBucketIAMMember{}
	Google := &providerv1alpha1.Google{}
	EtcdV3 := &backendv1alpha1.EtcdV3{}

	for {
		e := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
		h := map[string][]string{"Metadata-Flavor": {"Google"}}
		if ret := checkURL(e, h, 200); ret == nil {
			break
		}
	}

	if err := r.Get(context.Background(), req.NamespacedName, GoogleStorageBucket); !errors.IsNotFound(err) {
		if util.IsBeingDeleted(GoogleStorageBucket) {
			if err := r.deletionReconcile(GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3); err != nil {
				return reconcile.Result{}, err
			}

			if err := terraform.WriteToFile([]byte("{}"), GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.terraformReconcile(GoogleStorageBucket); err != nil {
				return reconcile.Result{}, err
			}

			util.RemoveFinalizer(GoogleStorageBucket, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name)

			if err := r.Update(context.Background(), GoogleStorageBucket); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		if dependency_met, err := r.dependencyReconcile(GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3); err == nil {
			// Check if dependency is met else interate again
			if !dependency_met {
				// Set the data
				GoogleStorageBucket.Status.State = "Failure"
				GoogleStorageBucket.Status.Phase = "Dependency"
				// Update the CR with status success
				if err := r.Status().Update(context.Background(), GoogleStorageBucket); err != nil {
					return reconcile.Result{}, err
				}
				// Dependency not met, don't error but finish reconcile until next change
				return reconcile.Result{}, nil
			}

			if !reflect.DeepEqual("Dependency", GoogleStorageBucket.Status.Phase) {
				// Set the data
				GoogleStorageBucket.Status.State = "Success"
				GoogleStorageBucket.Status.Phase = "Dependency"
				// Update the CR with status ready
				if err := r.Status().Update(context.Background(), GoogleStorageBucket); err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			return reconcile.Result{}, err
		}

		// Add finalizer to the GoogleStorageBucket resource
		util.AddFinalizer(GoogleStorageBucket, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name)

		// Update the CR with finalizer
		if err := r.Update(context.Background(), GoogleStorageBucket); err != nil {
			return reconcile.Result{}, err
		}

		b, err := terraform.RenderModuleToTerraform(GoogleStorageBucket.Spec, strings.ToLower(GoogleStorageBucket.Kind)+"_"+strings.ToLower(GoogleStorageBucket.ObjectMeta.Name))
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.WriteToFile(b, GoogleStorageBucket.ObjectMeta.Namespace, "Module_"+GoogleStorageBucket.Kind+"_"+GoogleStorageBucket.ObjectMeta.Name)
		if err != nil {
			return reconcile.Result{}, err
		}

		if !reflect.DeepEqual("Success", GoogleStorageBucket.Status.State) {
			// Set the data
			GoogleStorageBucket.Status.State = "Success"
			GoogleStorageBucket.Status.Phase = "Output"

			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), GoogleStorageBucket); err != nil {
				return reconcile.Result{}, err
			}
		}

		if err := r.terraformReconcile(GoogleStorageBucket); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err := r.Get(context.Background(), req.NamespacedName, GoogleStorageBucketIAMMember); !errors.IsNotFound(err) {
		if util.IsBeingDeleted(GoogleStorageBucketIAMMember) {
			if err := r.deletionReconcile(GoogleStorageBucketIAMMember, GoogleStorageBucket, Google, EtcdV3); err != nil {
				return reconcile.Result{}, err
			}

			if err := terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.terraformReconcile(GoogleStorageBucketIAMMember); err != nil {
				return reconcile.Result{}, err
			}

			util.RemoveFinalizer(GoogleStorageBucketIAMMember, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name)

			if err := r.Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		if dependency_met, err := r.dependencyReconcile(GoogleStorageBucketIAMMember, GoogleStorageBucket, Google, EtcdV3); err == nil {
			// Check if dependency is met else interate again
			if !dependency_met {
				// Set the data
				GoogleStorageBucketIAMMember.Status.State = "Failure"
				GoogleStorageBucketIAMMember.Status.Phase = "Dependency"
				// Update the CR with status success
				if err := r.Status().Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
					return reconcile.Result{}, err
				}
				// Dependency not met, don't error but finish reconcile until next change
				return reconcile.Result{}, nil
			}

			if !reflect.DeepEqual("Dependency", GoogleStorageBucketIAMMember.Status.Phase) {
				// Set the data
				GoogleStorageBucketIAMMember.Status.State = "Success"
				GoogleStorageBucketIAMMember.Status.Phase = "Dependency"
				// Update the CR with status ready
				if err := r.Status().Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			return reconcile.Result{}, err
		}

		// Add finalizer to the GoogleStorageBucketIAMMember resource
		util.AddFinalizer(GoogleStorageBucketIAMMember, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name)

		// Update the CR with finalizer
		if err := r.Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
			return reconcile.Result{}, err
		}

		c, err := terraform.RenderModuleToTerraform(GoogleStorageBucketIAMMember.Spec, strings.ToLower(GoogleStorageBucketIAMMember.Kind)+"_"+strings.ToLower(GoogleStorageBucketIAMMember.ObjectMeta.Name))
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.WriteToFile(c, GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name)
		if err != nil {
			return reconcile.Result{}, err
		}

		if !reflect.DeepEqual("Success", GoogleStorageBucketIAMMember.Status.State) {
			// Set the data
			GoogleStorageBucketIAMMember.Status.State = "Success"
			GoogleStorageBucketIAMMember.Status.Phase = "Output"

			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
				return reconcile.Result{}, err
			}
		}

		if err := r.terraformReconcile(GoogleStorageBucketIAMMember); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileModule) SetupWithGoogleStorageBucket(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.GoogleStorageBucket{}).
		Watches(&source.Kind{Type: &modulev1alpha1.GoogleStorageBucket{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}

func (r *ReconcileModule) SetupWithGoogleStorageBucketIAMMember(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.GoogleStorageBucketIAMMember{}).
		Watches(&source.Kind{Type: &modulev1alpha1.GoogleStorageBucketIAMMember{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
