package provider

import (
	"os"
	"time"
	"context"
	"reflect"

	terraformv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
	terraform "github.com/krubot/terraform-operator/pkg/terraform"
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

var log = logf.Log.WithName("controller_provider")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Provider Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProvider{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("provider-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Provider
	err = c.Watch(&source.Kind{Type: &terraformv1alpha1.Provider{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Provider
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &terraformv1alpha1.Provider{},
	})
	if err != nil {
		return err
	}

	return nil
}

// Path of terraform workspace
var TFPATH = os.Getenv("TFPATH")

// Provider structure to render the file
type Provider struct {
	Provider map[string]interface{} `json:"provider"`
}

// blank assignment to verify that ReconcileProvider implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileProvider{}

// ReconcileProvider reconciles a Provider object
type ReconcileProvider struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Provider object and makes changes based on the state read
// and what is in the Provider.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileProvider) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Provider")

	// Fetch the Provider instance
	instance := &terraformv1alpha1.Provider{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	// Update CR with the AppStatus == Starting
	if !reflect.DeepEqual("Started", instance.Status) {
		// Set the data
		instance.Status = "Started"

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Provider")
			return reconcile.Result{}, err
		}
	}

	time.Sleep(2 * time.Second)

	b, err := terraform.RenderProviderToTerraform(instance.Spec, instance.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = terraform.WriteToFile(b, instance.ObjectMeta.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("Created", instance.Status) {
		// Set the data
		instance.Status = "Created"

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Provider")
			return reconcile.Result{}, err
		}
	}

	time.Sleep(2 * time.Second)

	err = terraform.TerraformInit()
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("initialized", instance.Status) {
		// Set the data
		instance.Status = "initialized"

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Provider")
			return reconcile.Result{}, err
		}
	}

	time.Sleep(2 * time.Second)

	err = terraform.TerraformValidate()
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("validated", instance.Status) {
		// Set the data
		instance.Status = "validated"

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Provider")
			return reconcile.Result{}, err
		}
	}

	time.Sleep(2 * time.Second)

	err = terraform.TerraformPlan()
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update CR with the AppStatus == Created
	if !reflect.DeepEqual("planned", instance.Status) {
		// Set the data
		instance.Status = "planned"

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Project Status for the Provider")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
