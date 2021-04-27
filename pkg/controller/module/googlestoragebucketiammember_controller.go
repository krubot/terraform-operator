package controllers

import (
	"context"
	"fmt"
	"os"
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

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

// ReconcileProvider reconciles a Backend object
type ReconcileGoogleStorageBucketIAMMember struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ReconcileGoogleStorageBucketIAMMember) deletionReconcileGoogleStorageBucketIAMMember(module *modulev1alpha1.GoogleStorageBucketIAMMember, finalizerInterfaces ...interface{}) error {
	for _, finalizerInterface := range finalizerInterfaces {
		for _, fin := range module.GetFinalizers() {
			instance_split_fin := strings.Split(fin, "_")
			switch finalizer := finalizerInterface.(type) {
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
			case *modulev1alpha1.GoogleStorageBucket:
				if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(module, fin)
						if err := r.Update(context.Background(), module); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GoogleStorageBucket dependency is not met for deletion")
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
			case *backendv1alpha1.GCS:
				if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "GCS" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(module, fin)
						if err := r.Update(context.Background(), module); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GCS dependency is not met for deletion")
					}
				}
			}
		}
	}
	return nil
}

func (r *ReconcileGoogleStorageBucketIAMMember) dependencyReconcileGoogleStorageBucketIAMMember(module *modulev1alpha1.GoogleStorageBucketIAMMember, depInterfaces ...interface{}) (bool, error) {
	// Set the initial depency state
	dependency_met := true
	// Over the list of dependencies
	for _, depInterface := range depInterfaces {
		for _, depModule := range module.Dep {
			switch dep := depInterface.(type) {
			case *modulev1alpha1.GoogleStorageBucket:
				if depModule.Kind == "Module" && depModule.Type == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
				if depModule.Kind == "Module" && depModule.Type == "GoogleStorageBucketIAMMember" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
				if depModule.Kind == "Provider" && depModule.Type == "Google" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
				if depModule.Kind == "Backend" && depModule.Type == "EtcdV3" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
			case *backendv1alpha1.GCS:
				if depModule.Kind == "Backend" && depModule.Type == "GCS" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depModule.Name, Namespace: module.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
	return dependency_met, nil
}

func (r *ReconcileGoogleStorageBucketIAMMember) terraformReconcileGoogleStorageBucketIAMMember(module *modulev1alpha1.GoogleStorageBucketIAMMember) error {
	if err := terraform.TerraformInit(module.ObjectMeta.Namespace, os.Getenv("USER_WORKDIR")); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Init"
		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return err
		}
		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
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

	if err := terraform.TerraformNewWorkspace(module.ObjectMeta.Namespace, os.Getenv("USER_WORKDIR")); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Workspace"
		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return err
		}
		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
			return err
		}
		return nil
	}

	if err := terraform.TerraformSelectWorkspace(module.ObjectMeta.Namespace, os.Getenv("USER_WORKDIR")); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Workspace"
		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return err
		}
		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
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

	if err := terraform.TerraformApply(module.ObjectMeta.Namespace, os.Getenv("USER_WORKDIR")); err != nil {
		// Set the data
		module.Status.State = "Failure"
		module.Status.Phase = "Apply"
		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), module); err != nil {
			return err
		}
		if err := terraform.RemoveFile(module.ObjectMeta.Namespace, module.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
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

	for i := 0; i < reflect.TypeOf(module.Output).NumField(); i++ {
		output, err := terraform.TerraformOutput(module.ObjectMeta.Namespace, os.Getenv("USER_WORKDIR"), strings.ToLower(module.Kind)+"_"+strings.ToLower(module.ObjectMeta.Name)+"_"+reflect.TypeOf(module.Output).Field(i).Tag.Get("json"))
		if err != nil {
			return nil
		}

		t := reflect.ValueOf(&module.Output).Elem()
		val := t.FieldByName(reflect.TypeOf(module.Output).Field(i).Name)

		if val.CanSet() {
			val.SetString(output)
		}
	}

	// Update the CR with status ready
	if err := r.Update(context.Background(), module); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileGoogleStorageBucketIAMMember) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	GoogleStorageBucket := &modulev1alpha1.GoogleStorageBucket{}
	GoogleStorageBucketIAMMember := &modulev1alpha1.GoogleStorageBucketIAMMember{}
	Google := &providerv1alpha1.Google{}
	EtcdV3 := &backendv1alpha1.EtcdV3{}
	GCS := &backendv1alpha1.GCS{}

	for {
		e := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
		h := map[string][]string{"Metadata-Flavor": {"Google"}}
		if ret := util.CheckURL(e, h, 200); ret == nil {
			break
		}
	}

	if err := r.Get(context.Background(), req.NamespacedName, GoogleStorageBucketIAMMember); !errors.IsNotFound(err) {
		if util.IsBeingDeleted(GoogleStorageBucketIAMMember) {
			if err := r.deletionReconcileGoogleStorageBucketIAMMember(GoogleStorageBucketIAMMember, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3, GCS); err != nil {
				return reconcile.Result{}, err
			}

			d, err := terraform.RenderOutputToTerraform(GoogleStorageBucketIAMMember.Output, strings.ToLower(GoogleStorageBucketIAMMember.Kind)+"_"+strings.ToLower(GoogleStorageBucketIAMMember.ObjectMeta.Name))
			if err != nil {
				return reconcile.Result{}, err
			}

			for i, _ := range d {
				err = terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Output_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name+"_"+fmt.Sprint(i), os.Getenv("USER_WORKDIR"))
				if err != nil {
					return reconcile.Result{}, err
				}
			}

			if err := terraform.WriteToFile([]byte("{}"), GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.terraformReconcileGoogleStorageBucketIAMMember(GoogleStorageBucketIAMMember); err != nil {
				return reconcile.Result{}, err
			}

			util.RemoveFinalizer(GoogleStorageBucketIAMMember, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name)

			if err := r.Update(context.Background(), GoogleStorageBucketIAMMember); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		if dependency_met, err := r.dependencyReconcileGoogleStorageBucketIAMMember(GoogleStorageBucketIAMMember, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3, GCS); err == nil {
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
				return reconcile.Result{}, errors.NewBadRequest("GoogleStorageBucketIAMMember dependencies have not been met")
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

		err = terraform.WriteToFile(c, GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Module_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name, os.Getenv("USER_WORKDIR"))
		if err != nil {
			return reconcile.Result{}, err
		}

		d, err := terraform.RenderOutputToTerraform(GoogleStorageBucketIAMMember.Output, strings.ToLower(GoogleStorageBucketIAMMember.Kind)+"_"+strings.ToLower(GoogleStorageBucketIAMMember.ObjectMeta.Name))
		if err != nil {
			return reconcile.Result{}, err
		}

		for i, o := range d {
			err = terraform.WriteToFile(o, GoogleStorageBucketIAMMember.ObjectMeta.Namespace, "Output_"+GoogleStorageBucketIAMMember.Kind+"_"+GoogleStorageBucketIAMMember.ObjectMeta.Name+"_"+fmt.Sprint(i), os.Getenv("USER_WORKDIR"))
			if err != nil {
				return reconcile.Result{}, err
			}
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

		if err := r.terraformReconcileGoogleStorageBucketIAMMember(GoogleStorageBucketIAMMember); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGoogleStorageBucketIAMMember) SetupWithGoogleStorageBucketIAMMember(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.GoogleStorageBucketIAMMember{}).
		Watches(&source.Kind{Type: &modulev1alpha1.GoogleStorageBucketIAMMember{}}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
