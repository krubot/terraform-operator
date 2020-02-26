package terraform

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/terraform-svchost/auth"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/command/cliconfig"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/hashicorp/terraform/version"
	"github.com/mitchellh/cli"

	backendInit "github.com/hashicorp/terraform/backend/init"
	pluginDiscovery "github.com/hashicorp/terraform/plugin/discovery"
)

func credentialsSource(config *cliconfig.Config) (auth.CredentialsSource, error) {
	helperPlugins := pluginDiscovery.FindPlugins("credentials", globalPluginDirs())
	return config.CredentialsSource(helperPlugins)
}

func globalPluginDirs() []string {
	var ret []string
	// Look in ~/.terraform.d/plugins/ , or its equivalent on non-UNIX
	dir, err := cliconfig.ConfigDir()
	if err != nil {
		log.Printf("[ERROR] Error finding global config directory: %s", err)
	} else {
		machineDir := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
		ret = append(ret, filepath.Join(dir, "plugins"))
		ret = append(ret, filepath.Join(dir, "plugins", machineDir))
	}
	return ret
}

func TerraformNewWorkspace(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.WorkspaceNewCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace, namespace})
	log.Print("TerraformNewWorkspace exit code:", exitCode)
	return nil
}

func TerraformSelectWorkspace(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.WorkspaceSelectCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace, namespace})
	log.Print("TerraformSelectWorkspace exit code:", exitCode)
	return nil
}

func TerraformInit(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.InitCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace})
	log.Print("TerraformInit exit code:", exitCode)
	if exitCode != 0 {
		return errors.New("Terraform init returned a none zero exit code")
	}
	return nil
}

func TerraformValidate(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.ValidateCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace})
	log.Print("TerraformValidate exit code:", exitCode)
	if exitCode != 0 {
		return errors.New("Terraform validate returned a none zero exit code")
	}
	return nil
}

func TerraformPlan(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.PlanCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace})
	log.Print("TerraformPlan exit code:", exitCode)
	if exitCode != 0 {
		return errors.New("Terraform plane returned a none zero exit code")
	}
	return nil
}

func TerraformApply(namespace string) error {
	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	initCmd := &command.ApplyCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{"-auto-approve", namespace})
	log.Print("TerraformApply exit code:", exitCode)
	if exitCode != 0 {
		return errors.New("Terraform apply returned a none zero exit code")
	}
	return nil
}
