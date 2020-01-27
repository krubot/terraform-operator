package terraform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mitchellh/cli"
	"github.com/hashicorp/terraform/command"
	backendInit "github.com/hashicorp/terraform/backend/init"
)

func TerraformNewWorkspace(namespace string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("terraform", "workspace", "new", namespace)
	cmd.Dir = os.Getwd + "/" + namespace
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
	cmd.Dir = os.Getwd + "/" + namespace
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	fmt.Println("terraform init output:\n" + out.String())
	return nil
}

func TerraformInit(namespace string) error {
	backendInit.Init(nil)

  ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
  }

  initCmd := &command.InitCommand{
  	Meta: command.Meta{
  		Ui: ui,
  	},
  }

	exitStatus, err := initCmd.Run()
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
	cmd.Dir = os.Getwd + "/" + namespace
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
	cmd.Dir = os.Getwd + "/" + namespace
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
	cmd.Dir = os.Getwd + "/" + namespace
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
