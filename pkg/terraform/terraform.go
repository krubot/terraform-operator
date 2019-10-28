package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var TFPATH = os.Getenv("TFPATH")

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
	err := os.MkdirAll(TFPATH+"/"+namespace, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(TFPATH+"/"+namespace+"/"+name+".tf.json", b, 0755)
	if err != nil {
		return err
	}
	return nil
}

func TerraformNewWorkspace(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "workspace", "new", namespace)
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	fmt.Println("terraform init output:\n" + out.String())
	return nil
}

func TerraformSelectWorkspace(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "workspace", "select", namespace)
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	fmt.Println("terraform init output:\n" + out.String())
	return nil
}

func TerraformInit(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "init")
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("terraform init output:\n" + out.String())
	return nil
}

func TerraformValidate(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "validate")
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("terraform validate output:\n" + out.String())
	return nil
}

func TerraformPlan(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "plan")
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("terraform plan output:\n" + out.String())
	return nil
}

func TerraformApply(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = TFPATH + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("terraform apply output:\n" + out.String())
	return nil
}
