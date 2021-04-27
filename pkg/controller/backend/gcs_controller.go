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

// ReconcileGCS reconciles a Backend object
type ReconcileGCS struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ReconcileGCS) deletionReconcileGCS(backend *backendv1alpha1.GCS, finalizerInterfaces ...interface{}) error {
	for _, finalizerInterface := range finalizerInterfaces {
		for _, fin := range backend.GetFinalizers() {
			instance_split_fin := strings.Split(fin, "_")
			switch finalizer := finalizerInterface.(type) {
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
			case *backendv1alpha1.GCS:
				if instance_split_fin[0] == "Backend" && instance_split_fin[1] == "GCS" && instance_split_fin[2] != backend.ObjectMeta.Name {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(backend, fin)
						if err := r.Update(context.Background(), backend); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GCS dependency is not met for deletion")
					}
				}
			case *modulev1alpha1.GoogleStorageBucket:
				if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(backend, fin)
						if err := r.Update(context.Background(), backend); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GoogleStorageBucket dependency is not met for deletion")
					}
				}
			case *modulev1alpha1.GoogleStorageBucketIAMMember:
				if instance_split_fin[0] == "Module" && instance_split_fin[1] == "GoogleStorageBucketIAMMember" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
						util.RemoveFinalizer(backend, fin)
						if err := r.Update(context.Background(), backend); err != nil {
							return err
						}
					} else {
						return errors.NewBadRequest("GoogleStorageBucketIAMMember dependency is not met for deletion")
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
			}
		}
	}
	return nil
}

func (r *ReconcileGCS) dependencyReconcileGCS(backend *backendv1alpha1.GCS, depInterfaces ...interface{}) (bool, error) {
	// Set the initial depency state
	dependency_met := true
	// Over the list of dependencies
	for _, depInterface := range depInterfaces {
		for _, depBackend := range backend.Dep {
			switch dep := depInterface.(type) {
			case *modulev1alpha1.GoogleStorageBucket:
				if depBackend.Kind == "Module" && depBackend.Type == "GoogleStorageBucket" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
			case *modulev1alpha1.GoogleStorageBucketIAMMember:
				if depBackend.Kind == "Module" && depBackend.Type == "GoogleStorageBucketIAMMember" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
				if depBackend.Kind == "Provider" && depBackend.Type == "Google" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
				if depBackend.Kind == "Backend" && depBackend.Type == "EtcdV3" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
			case *backendv1alpha1.GCS:
				if depBackend.Kind == "Backend" && depBackend.Type == "GCS" {
					if err := r.Get(context.Background(), types.NamespacedName{Name: depBackend.Name, Namespace: backend.ObjectMeta.Namespace}, dep); err != nil {
						return false, nil
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
	return dependency_met, nil
}

func (r *ReconcileGCS) Reconcile(req ctrl.Request) (ctrl.Result, error) {
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

	if err := r.Get(context.Background(), req.NamespacedName, GCS); !errors.IsNotFound(err) {
		if util.IsBeingDeleted(GCS) {
			if err := r.deletionReconcileGCS(GCS, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3, GCS); err != nil {
				return reconcile.Result{}, err
			}

			if err := terraform.WriteToFile([]byte("{}"), GCS.ObjectMeta.Namespace, "Backend_"+GCS.Kind+"_"+GCS.ObjectMeta.Name, os.Getenv("USER_WORKDIR")); err != nil {
				return reconcile.Result{}, err
			}

			util.RemoveFinalizer(GCS, "Backend_"+GCS.Kind+"_"+GCS.ObjectMeta.Name)

			if err := r.Update(context.Background(), GCS); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		if dependency_met, err := r.dependencyReconcileGCS(GCS, GoogleStorageBucket, GoogleStorageBucketIAMMember, Google, EtcdV3, GCS); err == nil {
			// Check if dependency is met else interate again
			if !dependency_met {
				// Set the data
				GCS.Status.State = "Failure"
				GCS.Status.Phase = "Dependency"
				// Update the CR with status success
				if err := r.Status().Update(context.Background(), GCS); err != nil {
					return reconcile.Result{}, err
				}
				// Dependency not met, don't error but finish reconcile until next change
				return reconcile.Result{}, errors.NewBadRequest("GCS dependencies have not been met")
			}

			if !reflect.DeepEqual("Dependency", GCS.Status.Phase) {
				// Set the data
				GCS.Status.State = "Success"
				GCS.Status.Phase = "Dependency"
				// Update the CR with status ready

				if err := r.Status().Update(context.Background(), GCS); err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			return reconcile.Result{}, err
		}

		// Add finalizer to the module resource
		util.AddFinalizer(GCS, "Backend_"+GCS.Kind+"_"+GCS.ObjectMeta.Name)

		// Update the CR with finalizer
		if err := r.Update(context.Background(), GCS); err != nil {
			return reconcile.Result{}, err
		}

		b, err := terraform.RenderBackendToTerraform(GCS.Spec, strings.ToLower(GCS.Kind))
		if err != nil {
			return reconcile.Result{}, err
		}

		err = terraform.WriteToFile(b, GCS.ObjectMeta.Namespace, "Backend_"+GCS.Kind+"_"+GCS.ObjectMeta.Name, os.Getenv("USER_WORKDIR"))
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update CR with the AppStatus == Created
		if !reflect.DeepEqual("Ready", GCS.Status.State) {
			// Set the data
			GCS.Status.State = "Success"
			GCS.Status.Phase = "Output"

			// Update the CR with status success
			if err = r.Status().Update(context.Background(), GCS); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGCS) SetupWithGCS(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backendv1alpha1.GCS{}).
		Watches(&source.Kind{Type: &backendv1alpha1.GCS{}}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
