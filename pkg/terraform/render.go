package terraform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
)

type Terraform struct {
	Terraform Backend `json:"terraform"`
}

type Backend struct {
	Backend map[string]interface{} `json:"backend"`
}

type Provider struct {
	Provider map[string]interface{} `json:"provider"`
}

type Module struct {
	Module map[string]interface{} `json:"module"`
}

type Output struct {
	Output map[string]interface{} `json:"output"`
}

func RenderOutputToTerraform(instance interface{}, moduleName string) ([][]byte, error) {
	var byte_list [][]byte
	v := reflect.TypeOf(instance)

	for i := 0; i < v.NumField(); i++ {
		t := Output{
			Output: map[string]interface{}{
				v.Field(i).Tag.Get("json"): map[string]string{
					"value": "${module." + moduleName + "." + v.Field(i).Tag.Get("json") + "}",
				},
			},
		}

		b, err := json.MarshalIndent(t, "", "\t")
		if err != nil {
			return byte_list, err
		}
		byte_list = append(byte_list, b)
	}
	return byte_list, nil
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

// RenderBackendToTerraform takes an object, and attempts to construct the appropriate terraform json from it.
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

func WriteToFile(b []byte, namespace string, name string, path string) error {
	if err := os.MkdirAll(path+"/"+namespace, os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path+"/"+namespace+"/"+name+".tf.json", b, 0755); err != nil {
		return err
	}
	return nil
}

func RemoveFile(namespace string, name string, path string) error {
	if err := os.Remove(path + "/" + namespace + "/" + name + ".tf.json"); err != nil {
		return err
	}
	return nil
}
