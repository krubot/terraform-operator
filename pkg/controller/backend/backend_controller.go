package backend

import (
	"context"
	"os"
	"reflect"

	terraformv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
	terraform "github.com/krubot/terraform-operator/pkg/terraform"
	util "github.com/krubot/terraform-operator/pkg/util"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_backend")

// Add creates a new Backend Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackend{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("backend-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Backend
	err = c.Watch(&source.Kind{Type: &terraformv1alpha1.Backend{}}, &handler.EnqueueRequestForObject{}, util.ResourceGenerationOrFinalizerChangedPredicate{})
	if err != nil {
		return err
	}

	// Watch for changes to namespace resources
	err = c.Watch(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &terraformv1alpha1.Backend{},
	})
	if err != nil {
		return err
	}

	return nil
}

// Path of terraform workspace
var TFPATH = os.Getenv("TFPATH")

// Backend structure to render the file
type Backend struct {
	Backend map[string]interface{} `json:"backend"`
}

// blank assignment to verify that ReconcileBackend implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBackend{}

// ReconcileBackend reconciles a Backend object
type ReconcileBackend struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Backend object and makes changes based on the state read
// and what is in the Backend.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackend) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backend")

	// Fetch the Backend instance
	instance := &terraformv1alpha1.Backend{}
	
	err := r.client.Get(context.Background(), request.NamespacedName, instance)
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

	namespaceList, _ := listNamespaces(r.client)

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
		err = r.client.Status().Update(context.Background(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Backend")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func listNamespaces(c client.Client) (corev1.NamespaceList, error) {
	// Fetch the Namespace list instance
	backendNamespaceList := corev1.NamespaceList{}
	backendOpts := client.ListOptions{}

	// This is a hack, sometimes we can return nothing so we need to cycle till we get something
	// Fill free to tell me what I'm doing wrong here!
  for len(backendNamespaceList.Items) == 0 {
	  if err := c.List(context.Background(), &backendOpts, &backendNamespaceList); err != nil {
	  	return backendNamespaceList, err
	  }
  }

	return backendNamespaceList, nil
}
