package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	backendv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/backend/v1alpha1"
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

// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.my.domain,resources=cronjobs/status,verbs=get;update;patch

func (r *ReconcileBackend) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// Fetch the Backend instance
	instance := &backendv1alpha1.EtcdV3{}

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

	b, err := terraform.RenderBackendToTerraform(instance.Spec, instance.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	namespaceList, err := listNamespaces(r)
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, v := range namespaceList.Items {
		err = terraform.WriteToFile(b, v.Name, instance.ObjectMeta.Name)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("Ready", instance.Status) {
		// Set the data
		instance.Status = "Ready"

		// Update the CR
		err = r.Status().Update(context.Background(), instance)
		if err != nil {
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
		Watches(&source.Kind{Type: &corev1.Namespace{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
					return []reconcile.Request{
						// Trigger a reconcile on the etcdv3 backend update, please add more backend definitions as the api expands
						{NamespacedName: types.NamespacedName{
							Name:      "etcdv3",
							Namespace: "",
						}},
					}
				}),
			}).
		WithEventFilter(util.ResourceGenerationOrFinalizerChangedPredicate{}).
		Complete(r)
}
