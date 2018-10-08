
package terraform

import (
	"fmt"
	"bytes"
	"os/exec"
	"io/ioutil"
	"encoding/json"
)

const TFPATH = "/etc/infra/"

type Resource struct {
	Resource map[string]interface{} `json:"resource"`
}

// RenderToTerraform takes an object, and attempts to construct the appropriate terraform json from it.
func RenderToTerraform(i interface{}, resourceName, instanceName string) ([]byte, error) {
	r := Resource{
		Resource: map[string]interface{}{
			resourceName: map[string]interface{}{
				instanceName: i,
			},
		},
	}
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return b, err
	}
	return b, nil
}

func WriteToFile(b []byte, name string) error {
	err := ioutil.WriteFile(TFPATH+"/"+name+".tf", b, 0755)
	if err != nil {
		return err
	}
	return nil
}

func TerraformValidate() error {
	var out bytes.Buffer
	cmd := exec.Command("./run-terraform-validate.sh")
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("terraform run output:\n%q", out.String())
	return nil
}
