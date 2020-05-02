# New module type documentation

These docs will go through setting up a new module API extensions for the terraform operator. It'll go through the process of adding a new terraform module and integrating it with the operator.

## Adding a new module

All modules get added to the modules folder found [here](../modules). We divide each module into the cloud provider that uses it i.e. `gcp` modules go into that folder.

## Extending API struct

So each of the module variables is a struct we need to represent in the terraform operator spec. These struct mappings link with terraform json formatting so that we can represent them in `golang`. New structs are added to `apis/modules` folder in a file named `<module-name>_types.go` under the correct api version.

In this file we have at least 5 structs `<module>`, `<module>List`, `<module>Spec`, `DepSpec` and `StatusSpec` where module is its name.

```go
type Module struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModuleSpec `json:"spec,omitempty"`
	Dep    []DepSpec  `json:"dep,omitempty"`
	Status StatusSpec `json:"status,omitempty"`
}

type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Module `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Module{}, &ModuleList{})
}
```

Next the spec is defined as a mapping to the module variable resources. Let's take the following variables for example:

```terraform
variable "project" {
  description = "Bucket project id."
  type        = string
}

variable "entities" {
  description = "Entities list to add the IAM policies/bindings."
  type        = list(string)
}

variable "labels" {
  description = "Labels to be attached to the buckets."
  type        = map(string)
}

variable "lifecycle_rules" {
  description = "List of lifecycle rules to configure."
  type = set(object({
    action    = map(string)
    condition = map(string)
  }))
}
```

This would translate to the following `golang` struct output:

```go
type ModuleSpec struct {
  // +kubebuilder:validation:Enum={"/opt/modules/gcp/module/"}
  Source string `json:"source"`
  // Bucket project id.
  Project string `json:"project"`
  // Entities list to add the IAM policies/bindings
  Entities []string `json:"entities"`
  // Labels to be attached to the buckets
  Labels map[string]string `json:"labels"`
  // List of lifecycle rules to configure.
  LifecycleRules LifecycleRule `json:"lifecycle_rules"`
}

type LifecycleRule struct {
  Action map[string]string `json:"action"`
  Condition map[string]string `json:"condition"`
}
```

This would need to be applied to your module use case instead.

## Update the module controller

Now with the module api accessible from within the operator we need to update the controller logic to reconcile any custom resource changes. Since each controller needs to be aware of every type we'll need to update `backend_controller.go`, `module_controller.go` and `provider_controller.go`.

The changes needed for both `backend_controller.go` and `provider_controller.go` will be very similar since we don't need to get rapped up in the modules specific reconcile logic. Firstly `deletionReconcile` will need a new case as follows:

```go
  case *modulev1alpha1.APIType:
	if instance_split_fin[0] == "Module" && instance_split_fin[1] == "APIType" {
		if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: backend.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
			util.RemoveFinalizer(backend, fin)
			if err := r.Update(context.Background(), backend); err != nil {
				return err
			}
		} else {
			return errors.NewBadRequest("APIType dependency is not met for deletion")
		}
	}
```

and

```go
  case *modulev1alpha1.APIType:
	if instance_split_fin[0] == "Module" && instance_split_fin[1] == "APIType" {
		if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: provider.ObjectMeta.Namespace}, finalizer); errors.IsNotFound(err) {
			util.RemoveFinalizer(provider, fin)
			if err := r.Update(context.Background(), provider); err != nil {
				return err
			}
		} else {
			return errors.NewBadRequest("APIType dependency is not met for deletion")
		}
	}
```

Additionally we need to update the `dependencyReconcile` with the following too:

```go
  case *modulev1alpha1.APIType:
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
```

and

```go
  case *modulev1alpha1.APIType:
	if depProvider.Kind == "Module" {
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
```

Lastly here we need to update the reconcile function itself with the additional API it will be making changes too.

```go
APIType := &modulev1alpha1.APIType{}
...

...
if util.IsBeingDeleted(EtcdV3) {
  if err := r.deletionReconcile(EtcdV3, APIType, Google); err != nil {
  	return reconcile.Result{}, err
  }
...

...
if dependency_met, err := r.dependencyReconcile(EtcdV3, APIType, Google); err == nil {
		if !dependency_met {
...
```

and

```go
APIType := &modulev1alpha1.APIType{}
...

...
if util.IsBeingDeleted(Google) {
  if err := r.deletionReconcile(Google, APIType, EtcdV3); err != nil {
  	return reconcile.Result{}, err
  }
...

...
if dependency_met, err := r.dependencyReconcile(Google, APIType, EtcdV3); err == nil {
		if !dependency_met {
...
```

Finally we need to update the `module_controller.go` file with changes for reconciling this new API struct. Similar to above we need sections for `APIType` here however this team we'll need to duplicate the lower case again:

```go
...
case *modulev1alpha1.APIType:
	for _, fin := range module.GetFinalizers() {
		instance_split_fin := strings.Split(fin, "_")
		switch finalizer := finalizerInterface.(type) {
...
		case *modulev1alpha1.APIType:
			if instance_split_fin[0] == "Module" && instance_split_fin[1] == "APIType" && instance_split_fin[2] != module.ObjectMeta.Name {
				if err := r.Get(context.Background(), types.NamespacedName{Name: instance_split_fin[2], Namespace: module.ObjectMeta.Namespace},finalizer); errors.IsNotFound(err) {
					util.RemoveFinalizer(module, fin)
					if err := r.Update(context.Background(), module); err != nil {
						return err
					}
				} else {
					return errors.NewBadRequest("APIType dependency is not met for deletion")
				}
			}
```

and

```go
...
case *modulev1alpha1.APIType:
	for _, depModule := range module.Dep {
		switch dep := depInterface.(type) {
...
		case *modulev1alpha1.APIType:
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
```

Additionally, we need to add a terraform reconcile section to this part since we need to apply these changes when in the correct state.

```go
...
case *modulev1alpha1.APIType:
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
...
```

Lastly, in this part we need sections for the `Reconcile` and `ControllerManager`. These will do all the hard work got get the changes into the correct state so we need to be careful here. The change however just need us the replicate the previous section with the new module struct:

```go
APIType := &modulev1alpha1.APIType{}
...
if err := r.Get(context.Background(), req.NamespacedName, APIType); !errors.IsNotFound(err) {
	if util.IsBeingDeleted(APIType) {
		if err := r.deletionReconcile(APIType, GoogleStorageBucket, Google, EtcdV3); err != nil {
			return reconcile.Result{}, err
		}
...
```

and

```go
func (r *ReconcileModule) SetupWithAPIType(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.APIType{}).
		Watches(&source.Kind{Type: &modulev1alpha1.APIType{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(util.ResourceGenerationChangedPredicate{}).
		Complete(r)
}
```

Then we'll need to add controller logic to be called in the `main.go` file found in the `cmd/manager/` folder:

```go
...
if err = (&controllers.ReconcileModule{
	Client: mgr.GetClient(),
	Log:    ctrl.Log.WithName("controllers").WithName("Module"),
	Scheme: mgr.GetScheme(),
}).SetupWithAPIType(mgr); err != nil {
	setupLog.Error(err, "unable to create controller", "controller", "Module")
	os.Exit(1)
}
```
