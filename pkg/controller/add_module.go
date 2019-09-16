package controller

import (
	"github.com/krubot/terraform-operator/pkg/controller/module"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, module.Add)
}
