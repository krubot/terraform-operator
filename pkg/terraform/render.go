package terraform

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Provider struct {
	Provider map[string]interface{} `json:"provider"`
}

type Module struct {
	Module map[string]interface{} `json:"module"`
}

type Terraform struct {
	Terraform Backend `json:"terraform"`
}

type Backend struct {
	Backend map[string]interface{} `json:"backend"`
}

// RenderProviderToTerraform takes an object, and attempts to construct the appropriate terraform json from it.
func RenderProviderToTerraform(instance interface{}, providerName string) ([]byte, error) {
	r := Provider{
		Provider: map[string]interface{}{
			providerName: instance,
		},
	}
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return b, err
	}
	return b, nil
}

// RenderModuleToTerraform takes an object, and attempts to construct the appropriate terraform json from it.
func RenderModuleToTerraform(instance interface{}, moduleName string) ([]byte, error) {
	r := Module{
		Module: map[string]interface{}{
			moduleName: instance,
		},
	}
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return b, err
	}
	return b, nil
}

// RenderModuleToTerraform takes an object, and attempts to construct the appropriate terraform json from it.
func RenderBackendToTerraform(instance interface{}, backendName string) ([]byte, error) {
	r := Terraform{
		Backend{
			Backend: map[string]interface{}{
				backendName: instance,
			},
		},
	}
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return b, err
	}
	return b, nil
}

func WriteToFile(b []byte, namespace string, name string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(currentDir+"/"+namespace, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(currentDir+"/"+namespace+"/"+name+".tf.json", b, 0755)
	if err != nil {
		return err
	}
	return nil
}
