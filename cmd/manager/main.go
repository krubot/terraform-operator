/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/krubot/terraform-operator/pkg/admission"
	"github.com/krubot/terraform-operator/pkg/version"

	backendv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/backend/v1alpha1"
	modulev1alpha1 "github.com/krubot/terraform-operator/pkg/apis/module/v1alpha1"
	providerv1alpha1 "github.com/krubot/terraform-operator/pkg/apis/provider/v1alpha1"
	backendcontroller "github.com/krubot/terraform-operator/pkg/controller/backend"
	modulecontroller "github.com/krubot/terraform-operator/pkg/controller/module"
	providercontroller "github.com/krubot/terraform-operator/pkg/controller/provider"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = providerv1alpha1.AddToScheme(scheme)

	_ = backendv1alpha1.AddToScheme(scheme)

	_ = modulev1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var webhookAddr int
	var enableLeaderElection bool
	var healthAddr string

	flag.StringVar(&metricsAddr, "metrics-addr", ":8383", "The address the metric endpoint binds to.")
	flag.IntVar(&webhookAddr, "webhook-addr", 9443, "The address the webhook endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&healthAddr, "health-addr", ":9440", "The address the health endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	setupLog.Info("Version of terraform-operator", "version", version.Version)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: healthAddr,
		ReadinessEndpointName:  "/readyz",
		LivenessEndpointName:   "/healthz",
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "terraform-operator-leader-election",
		Port:                   webhookAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	mgr.AddHealthzCheck("Liveness", healthz.Ping)

	setupLog.Info("controller reconcile")

	if err = (&backendcontroller.ReconcileEtcdV3{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Backend").WithName("EtcdV3"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Backend")
		os.Exit(1)
	}

	if err = (&providercontroller.ReconcileGoogle{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Provider").WithName("Google"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Provider")
		os.Exit(1)
	}

	if err = (&modulecontroller.ReconcileGoogleStorageBucket{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Module").WithName("GoogleStorageBucket"),
		Scheme: mgr.GetScheme(),
	}).SetupWithGoogleStorageBucket(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Module")
		os.Exit(1)
	}

	if err = (&modulecontroller.ReconcileGoogleStorageBucketIAMMember{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Module").WithName("GoogleStorageBucketIAMMember"),
		Scheme: mgr.GetScheme(),
	}).SetupWithGoogleStorageBucketIAMMember(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Module")
		os.Exit(1)
	}

	// Setup webhooks
	setupLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	setupLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/validate-terraform", &webhook.Admission{Handler: &admission.TerraformValidator{Client: mgr.GetClient()}})

	mgr.AddReadyzCheck("Readiness", healthz.Ping)

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
