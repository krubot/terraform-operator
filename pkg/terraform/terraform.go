package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/logutils"
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

var UI cli.Ui

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

func logOutput() io.Writer {
	levels := []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	minLevel := os.Getenv("TF_LOG")

	// default log writer is null device.
	writer := ioutil.Discard
	if minLevel != "" {
		writer = os.Stderr
	}

	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}

	return filter
}

func TerraformNewWorkspace(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace})
	log.Print("TerraformNewWorkspace exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	return nil
}

func TerraformSelectWorkspace(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{namespace})
	log.Print("TerraformSelectWorkspace exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	return nil
}

func TerraformInit(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{})
	log.Print("TerraformInit exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return errors.New("Terraform init returned a none zero exit code")
	}
	return nil
}

func TerraformValidate(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{})
	log.Print("TerraformValidate exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return errors.New("Terraform validate returned a none zero exit code")
	}
	return nil
}

func TerraformPlan(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{"-out=" + namespace + ".tfplan"})
	log.Print("TerraformPlan exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return errors.New("Terraform plane returned a none zero exit code")
	}
	return nil
}

func TerraformShow(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	if err := os.MkdirAll("/opt/plan", os.ModePerm); err != nil {
		return err
	}

	fh, _ := os.OpenFile("/opt/plan/"+namespace+".tfplan.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.FileMode(0755))

	initCmd := &command.ShowCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Writer:      fh,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{"-json", namespace + ".tfplan"})
	log.Print("TerraformPlan exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return errors.New("Terraform plane returned a none zero exit code")
	}
	return nil
}

func TerraformApply(namespace string, path string) error {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

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
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}

	exitCode := initCmd.Run([]string{"-auto-approve"})
	log.Print("TerraformApply exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return errors.New("Terraform apply returned a none zero exit code")
	}
	return nil
}

func TerraformOutput(namespace string, path string, name string) (string, error) {
	if err := os.Chdir(filepath.Join(path, namespace)); err != nil {
		return "", errors.New("Couldn't change directory path")
	}

	log.SetOutput(logOutput())

	config, _ := cliconfig.LoadConfig()
	credsSrc, err := credentialsSource(config)
	if err != nil {
		log.Printf("[WARN] Cannot initialize remote host credentials manager: %s", err)
	}

	services := disco.NewWithCredentialsSource(credsSrc)
	services.SetUserAgent(httpclient.TerraformUserAgent(version.String()))

	backendInit.Init(services)

	var tplStdOut bytes.Buffer

	initCmd := &command.OutputCommand{
		Meta: command.Meta{
			Color:               false,
			RunningInAutomation: true,
			PluginCacheDir:      config.PluginCacheDir,
			Ui: &cli.BasicUi{
				Writer: &tplStdOut,
			},
		},
	}

	exitCode := initCmd.Run([]string{name})
	log.Print("TerraformApply exit code:", exitCode)
	if err := os.Chdir(path); err != nil {
		return "", errors.New("Couldn't change directory path")
	}
	if exitCode != 0 {
		return "", errors.New("Terraform apply returned a none zero exit code")
	}
	return strings.ReplaceAll(tplStdOut.String(), " ", ""), nil
}
