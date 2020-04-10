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
type ReconcileProvider struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

func (r *ReconcileProvider) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// Fetch the terraform instances
	provider := &providerv1alpha1.Google{}

	if err := r.Get(context.Background(), req.NamespacedName, provider); err != nil {
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

	if util.IsBeingDeleted(provider) {
		finalizersProvider := &providerv1alpha1.Google{}
		finalizersModule := &modulev1alpha1.GCS{}
		finalizersBackend := &backendv1alpha1.EtcdV3{}

		for _, fin := range provider.GetFinalizers() {
			instance_split_fin := strings.Split(fin, "_")

			if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" && instance_split_fin[2] != provider.ObjectMeta.Name {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersProvider); errors.IsNotFound(err) {
					util.RemoveFinalizer(provider, fin)

					if err := r.Update(context.Background(), provider); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("provider dependency is not met for deletion")
				}
			}

			if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersBackend); errors.IsNotFound(err) {
					util.RemoveFinalizer(provider, fin)

					if err := r.Update(context.Background(), provider); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("backend dependency is not met for deletion")
				}
			}

			if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GCS" {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: req.NamespacedName.Namespace}, finalizersModule); errors.IsNotFound(err) {
					util.RemoveFinalizer(provider, fin)

					if err := r.Update(context.Background(), provider); err != nil {
						return reconcile.Result{}, err
					}
				} else {
					return reconcile.Result{}, errors.NewBadRequest("module dependency is not met for deletion")
				}
			}
		}

		if err := terraform.WriteToFile([]byte("{}"), provider.ObjectMeta.Namespace, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		util.RemoveFinalizer(provider, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)

		if err := r.Update(context.Background(), provider); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// Set the initial depency state
	dependency_met := true

	// List of dependencies to add resource finalisers too
	for _, dep := range provider.Dep {
		depProvider := &providerv1alpha1.Google{}
		depModule := &modulev1alpha1.GCS{}
		depBackend := &backendv1alpha1.EtcdV3{}

		if dep.Kind == "Backend" {
			if err := r.Get(context.Background(), types.NamespacedName{Name: dep.Name, Namespace: req.NamespacedName.Namespace}, depBackend); err != nil {
				return reconcile.Result{}, err
			}

			if depBackend.Status.State == "Success" {
				// Add finalizer to the module resource
				util.AddFinalizer(depBackend, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)

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
				util.AddFinalizer(depModule, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)

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
				util.AddFinalizer(depProvider, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)

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
		provider.Status.State = "Failure"
		provider.Status.Phase = "Dependency"

		// Update the CR with status success
		if err := r.Status().Update(context.Background(), provider); err != nil {
			return reconcile.Result{}, err
		}

		// Dependency not met, don't error but finish reconcile until next change
		return reconcile.Result{}, nil
	}

	if !reflect.DeepEqual("Dependency", provider.Status.Phase) {
		// Set the data
		provider.Status.State = "Success"
		provider.Status.Phase = "Dependency"

		// Update the CR with status ready
		if err := r.Status().Update(context.Background(), provider); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Add finalizer to the module resource
	util.AddFinalizer(provider, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)

	// Update the CR with finalizer
	if err := r.Update(context.Background(), provider); err != nil {
		return reconcile.Result{}, err
	}

	b, err := terraform.RenderProviderToTerraform(provider.Spec, strings.ToLower(provider.Kind))
	if err != nil {
		return reconcile.Result{}, err
	}

	err = terraform.WriteToFile(b, provider.ObjectMeta.Namespace, "Provider_"+provider.Kind+"_"+provider.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("Ready", provider.Status.State) {
		// Set the data
		provider.Status.State = "Success"
		provider.Status.Phase = "Output"

		// Update the CR
		if err = r.Status().Update(context.Background(), provider); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileProvider) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&providerv1alpha1.Google{}).
		Watches(&source.Kind{Type: &providerv1alpha1.Google{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
