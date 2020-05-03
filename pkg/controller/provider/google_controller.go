package controllers

import (
	"context"
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

// ReconcileGoogle reconciles a Backend object
type ReconcileGoogle struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ReconcileGoogle) deletionReconcileGoogle(provider *providerv1alpha1.Google, finalizerInterfaces ...interface{}) error {
	for _, finalizerInterface := range finalizerInterfaces {
		for _, fin := range provider.GetFinalizers() {
			instance_split_fin := strings.Split(fin, "_")
			switch finalizer := finalizerInterface.(type) {
			case *providerv1alpha1.Google:
				if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" && instance_split_fin[2] != provider.ObjectMeta.Name {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: provider.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(provider, fin)
						if err := r.Update(context.Background(), provider); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("Google dependency is not met for deletion")
					}
				}
			case *modulev1alpha1.GoogleStorageBucket:
				if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: provider.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(provider, fin)
						if err := r.Update(context.Background(), provider); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GoogleStorageBucket dependency is not met for deletion")
					}
				}
			case *modulev1alpha1.GoogleStorageBucketIAMMember:
				if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucketIAMMember" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: provider.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(provider, fin)
						if err := r.Update(context.Background(), provider); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GoogleStorageBucketIAMMember dependency is not met for deletion")
					}
				}
			case *backendv1alpha1.EtcdV3:
				if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: provider.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(provider, fin)
						if err := r.Update(context.Background(), provider); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("EtcdV3 dependency is not met for deletion")
					}
				}
			}
		}
	}
	return nil
}

func (r *ReconcileGoogle) dependencyReconcileGoogle(provider *providerv1alpha1.Google, depInterfaces ...interface{}) (bool, error) {
	// Set the initial depency state
	dependency_met := true
	// Over the list of dependencies
	for _, depInterface := range depInterfaces {
		for _, depProvider := range provider.Dep {
			switch dep := depInterface.(type) {
			case *modulev1alpha1.GoogleStorageBucket:
				if depProvider.Kind == "Module" && depProvider.Type == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depProvider.Name, Namespace: provider.ObjectMeta.Namespace}, dep); err != nil {
						return dependency_met, err
					}
					if dep.Status.State == "Success" {
						// Add finalizer to the GoogleStorageBucket. resource
						util.AddFinalizer(dep, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)
						// Update the CR with finalizer
						if err := r.Update(context.Background(), dep); err != nil {
							return dependency_met, err
						}
					} else {
						dependency_met = false
					}
				}
			case *modulev1alpha1.GoogleStorageBucketIAMMember:
				if depProvider.Kind == "Module" && depProvider.Type == "GoogleStorageBucketIAMMember" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depProvider.Name, Namespace: provider.ObjectMeta.Namespace}, dep); err != nil {
						return dependency_met, err
					}
					if dep.Status.State == "Success" {
						// Add finalizer to the GoogleStorageBucket. resource
						util.AddFinalizer(dep, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)
						// Update the CR with finalizer
						if err := r.Update(context.Background(), dep); err != nil {
							return dependency_met, err
						}
					} else {
						dependency_met = false
					}
				}
			case *providerv1alpha1.Google:
				if depProvider.Kind == "Provider" && depProvider.Type == "Google" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depProvider.Name, Namespace: provider.ObjectMeta.Namespace}, dep); err != nil {
						return dependency_met, err
					}
					if dep.Status.State == "Success" {
						// Add finalizer to the GoogleStorageBucket resource
						util.AddFinalizer(dep, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)
						// Update the CR with finalizer
						if err := r.Update(context.Background(), dep); err != nil {
							return dependency_met, err
						}
					} else {
						dependency_met = false
					}
				}
			case *backendv1alpha1.EtcdV3:
				if depProvider.Kind == "Backend" && depProvider.Type == "EtcdV3" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depProvider.Name, Namespace: provider.ObjectMeta.Namespace}, dep); err != nil {
						return dependency_met, err
					}
					if dep.Status.State == "Success" {
						// Add finalizer to the GoogleStorageBucket resource
						util.AddFinalizer(dep, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)
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

func (r *ReconcileGoogle) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	GoogleStorageBucket := &modulev1alpha1.GoogleStorageBucket{}
	GoogleStorageBucketIAMMember := &modulev1alpha1.GoogleStorageBucketIAMMember{}
	Google := &providerv1alpha1.Google{}
	EtcdV3 := &backendv1alpha1.EtcdV3{}

	for {
		e := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
		h := map[string][]string{"Metadata-Flavor": {"Google"}}
		if ret := util.CheckURL(e, h, 200); ret == nil {
			break
		}
	}

	if err := r.Get(context.Background(), req.NamespacedName, Google); !errors.IsNotFound(err) {
		if util.IsBeingDeleted(Google) {
			if err := r.deletionReconcileGoogle(Google, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3); err != nil {
				return reconcile.Result{}, err
			}

			if err := terraform.WriteToFile([]byte("{}"), Google.ObjectMeta.Namespace, "Provider_"+Google.Kind+"_"+Google.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
				return reconcile.Result{}, err
			}

			util.RemoveFinalizer(Google, "Provider_"+Google.Kind+"_"+Google.ObjectMeta.Name)

			if err := r.Update(context.Background(), Google); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		if dependency_met, err := r.dependencyReconcileGoogle(Google, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3); err == nil {
			// Check if dependency is met else interate again
			if !dependency_met {
				// Set the data
				Google.Status.State = "Failure"
				Google.Status.Phase = "Dependency"
				// Update the CR with status success
				if err := r.Status().Update(context.Background(), Google); err != nil {
					return reconcile.Result{}, err
				}
				// Dependency not met, don't error but finish reconcile until next change
				return reconcile.Result{}, nil
			}

			if !reflect.DeepEqual("Dependency", Google.Status.Phase) {
				// Set the data
				Google.Status.State = "Success"
				Google.Status.Phase = "Dependency"
				// Update the CR with status ready
				if err := r.Status().Update(context.Background(), Google); err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			return reconcile.Result{}, err
		}

		// Add finalizer to the module resource
		util.AddFinalizer(Google, "Provider_"+Google.Kind+"_"+Google.ObjectMeta.Name)

		// Update the CR with finalizer
		if err := r.Update(context.Background(), Google); err != nil {
			return reconcile.Result{}, err
		}

		b, err := terraform.RenderProviderToTerraform(Google.Spec, strings.ToLower(Google.Kind))
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.WriteToFile(b, Google.ObjectMeta.Namespace, "Provider_"+Google.Kind+"_"+Google.ObjectMeta.Name, os.Getenv("USER_WORKDIR"))
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update CR with the AppStatus == Created
		if !reflect.DeepEqual("Ready", Google.Status.State) {
			// Set the data
			Google.Status.State = "Success"
			Google.Status.Phase = "Output"

			// Update the CR
			if err = r.Status().Update(context.Background(), Google); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGoogle) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&providerv1alpha1.Google{}).
		Watches(&source.Kind{Type: &providerv1alpha1.Google{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
