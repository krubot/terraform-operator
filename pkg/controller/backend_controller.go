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

// ReconcileBackend reconciles a Backend object
type ReconcileBackend struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ReconcileBackend) deletionReconcile(backendInterface interface{}, finalizerInterfaces ...interface{}) error {
	for _, finalizerInterface := range finalizerInterfaces {
		switch backend := backendInterface.(type) {
		case *backendv1alpha1.EtcdV3:
			for _, fin := range backend.GetFinalizers() {
				instance_split_fin := strings.Split(fin, "_")
				switch finalizer := finalizerInterface.(type) {
				case *modulev1alpha1.GCS:
					if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GCS" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(backend, fin)
							if err := r.Update(context.Background(), backend); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("GCS dependency is not met for deletion")
						}
					}
				case *providerv1alpha1.Google:
					if instance_split_fin[0] == "Provider" && instance_split_fin[1] == "Google" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(backend, fin)
							if err := r.Update(context.Background(), backend); err != nil {
								return err
							}
						} else {
							return errors.NewBadRequest("Google dependency is not met for deletion")
						}
					}
				case *backendv1alpha1.EtcdV3:
					if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "EtcdV3" && instance_split_fin[2] != backend.ObjectMeta.Name {
						if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
							util.RemoveFinalizer(backend, fin)
							if err := r.Update(context.Background(), backend); err != nil {
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

func (r *ReconcileBackend) dependencyReconcile(backendInterface interface{}, depInterfaces ...interface{}) (bool, error) {
	// Set the initial depency state
	dependency_met := true
	// Over the list of dependencies
	for _, depInterface := range depInterfaces {
		switch backend := backendInterface.(type) {
		case *backendv1alpha1.EtcdV3:
			for _, depBackend := range backend.Dep {
				switch dep := depInterface.(type) {
				case *modulev1alpha1.GCS:
					if depBackend.Kind == "Module" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							// Add finalizer to the GoogleStorageBucket. resource
							util.AddFinalizer(dep, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *providerv1alpha1.Google:
					if depBackend.Kind == "Provider" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							// Add finalizer to the GoogleStorageBucket resource
							util.AddFinalizer(dep, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)
							// Update the CR with finalizer
							if err := r.Update(context.Background(), dep); err != nil {
								return dependency_met, err
							}
						} else {
							dependency_met = false
						}
					}
				case *backendv1alpha1.EtcdV3:
					if depBackend.Kind == "Backend" {
						if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
							return dependency_met, err
						}
						if dep.Status.State == "Success" {
							// Add finalizer to the GoogleStorageBucket resource
							util.AddFinalizer(dep, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)
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

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

func (r *ReconcileBackend) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	module := &modulev1alpha1.GCS{}
	provider := &providerv1alpha1.Google{}
	backend := &backendv1alpha1.EtcdV3{}

	if err := r.Get(context.Background(), req.NamespacedName, backend); err != nil {
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

	if util.IsBeingDeleted(backend) {
		if err := r.deletionReconcile(backend, module, provider); err != nil {
			return reconcile.Result{}, err
		}

		if err := terraform.WriteToFile([]byte("{}"), backend.ObjectMeta.Namespace, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name); err != nil {
			return reconcile.Result{}, err
		}

		util.RemoveFinalizer(backend, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)

		if err := r.Update(context.Background(), backend); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if dependency_met, err := r.dependencyReconcile(backend, provider, module); err == nil {
		// Check if dependency is met else interate again
		if !dependency_met {
			// Set the data
			backend.Status.State = "Failure"
			backend.Status.Phase = "Dependency"
			// Update the CR with status success
			if err := r.Status().Update(context.Background(), backend); err != nil {
				return reconcile.Result{}, err
			}
			// Dependency not met, don't error but finish reconcile until next change
			return reconcile.Result{}, nil
		}

		if !reflect.DeepEqual("Dependency", backend.Status.Phase) {
			// Set the data
			backend.Status.State = "Success"
			backend.Status.Phase = "Dependency"
			// Update the CR with status ready
			if err := r.Status().Update(context.Background(), backend); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		return reconcile.Result{}, err
	}

	// Add finalizer to the module resource
	util.AddFinalizer(backend, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)

	// Update the CR with finalizer
	if err := r.Update(context.Background(), backend); err != nil {
		return reconcile.Result{}, err
	}

	b, err := terraform.RenderBackendToTerraform(backend.Spec, strings.ToLower(backend.Kind))
	if err != nil {
		return reconcile.Result{}, err
	}

	err = terraform.WriteToFile(b, backend.ObjectMeta.Namespace, "Backend_"+backend.Kind+"_"+backend.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("Ready", backend.Status.State) {
		// Set the data
		backend.Status.State = "Success"
		backend.Status.Phase = "Output"

		// Update the CR with status success
		if err = r.Status().Update(context.Background(), backend); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileBackend) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backendv1alpha1.EtcdV3{}).
		Watches(&source.Kind{Type: &backendv1alpha1.EtcdV3{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
