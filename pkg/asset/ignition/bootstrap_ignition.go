package ignition

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/ignition/v2/config/util"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	agentManifests "github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/ignition"
	ignasset "github.com/openshift/installer/pkg/asset/ignition"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"

	"github.com/openshift/appliance/pkg/asset/config"
	"github.com/openshift/appliance/pkg/asset/manifests"
	"github.com/openshift/appliance/pkg/consts"
	"github.com/openshift/appliance/pkg/templates"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/sirupsen/logrus"
)

const (
	bootstrapRegistryDataPath = "/mnt/agentdata/oc-mirror/bootstrap"
	registriesConfFilePath    = "/etc/containers/registries.conf"
	manifestPath              = "/etc/assisted/manifests"
	corePassOverrideFilePath  = "/etc/assisted/appliance-override-password-set"
	extraManifestPath         = "/etc/assisted/extra-manifests"
)

var (
	bootstrapServices = []string{
		"start-local-registry.service",
		"assisted-service.service",
		"create-cluster-and-infraenv.service",
		"pre-install.service",
		"pre-install-node-zero.service",
		"update-hosts.service",
	}

	bootstrapScripts = []string{
		"setup-local-registry.sh",
		"set-env-files.sh",
		"pre-install.sh",
		"pre-install-node-zero.sh",
		"release-image-download.sh",
		"release-image.sh",
		"update-hosts.sh",
		"create-virtual-device.sh",

		// TODO: remove (needed for using custom agent image)
		"get-container-images.sh",
	}
)

// BootstrapIgnition generates the bootstrap ignition file for the recovery ISO
type BootstrapIgnition struct {
	Config igntypes.Config
}

var _ asset.Asset = (*BootstrapIgnition)(nil)

// Name returns the human-friendly name of the asset.
func (i *BootstrapIgnition) Name() string {
	return "Bootstrap ignition"
}

// Dependencies returns dependencies used by the asset.
func (i *BootstrapIgnition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&config.EnvConfig{},
		&config.ApplianceConfig{},
		&password.KubeadminPassword{},
		&manifests.ClusterImageSet{},
		&agentManifests.ExtraManifests{},
		&InstallIgnition{},
	}
}

// Generate the base ISO.
func (i *BootstrapIgnition) Generate(dependencies asset.Parents) error {
	envConfig := &config.EnvConfig{}
	applianceConfig := &config.ApplianceConfig{}
	extraManifests := &agentManifests.ExtraManifests{}
	pwd := &password.KubeadminPassword{}
	installIgnition := &InstallIgnition{}
	dependencies.Get(envConfig, applianceConfig, extraManifests, pwd, installIgnition)

	i.Config = igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	if envConfig.DebugBootstrap {
		// Avoid machine reboot after bootstrap to debug install ignition
		bootstrapServices = append(bootstrapServices, "ironic-agent.service")
	}

	// Add services common for bootstrap and install
	if err := bootstrap.AddSystemdUnits(&i.Config, "services/common", nil, bootstrapServices); err != nil {
		return err
	}

	// Add services exclusive for bootstrap
	if err := bootstrap.AddSystemdUnits(&i.Config, "services/bootstrap", nil, bootstrapServices); err != nil {
		return err
	}

	// Fetch install ignition config
	installIgnitionConfig, err := json.Marshal(installIgnition.Config)
	if err != nil {
		return err
	}

	// Get base image path
	coreosImagePattern := fmt.Sprintf(consts.CoreosImagePattern, applianceConfig.GetCpuArchitecture())
	coreosImagePath := envConfig.FindInCache(coreosImagePattern)

	// Add bootstrap scripts to ignition
	templateData := templates.GetBootstrapIgnitionTemplateData(
		applianceConfig.Config.OcpRelease,
		bootstrapRegistryDataPath,
		string(installIgnitionConfig),
		coreosImagePath)
	for _, script := range bootstrapScripts {
		if err := bootstrap.AddStorageFiles(&i.Config,
			"/usr/local/bin/"+script,
			"scripts/bin/"+script+".template",
			templateData); err != nil {
			return err
		}
	}

	passwdUser := igntypes.PasswdUser{
		Name: "core",
	}
	// Add user 'core' password
	if applianceConfig.Config.UserCorePass != nil {
		passBytes, err := bcrypt.GenerateFromPassword([]byte(*applianceConfig.Config.UserCorePass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		pwdHash := string(passBytes)
		passwdUser.PasswordHash = &pwdHash

		// Add 'appliance-override-password-set' file
		// (needed as an indication that the appliance override the core pass)
		registryEnvFile := ignasset.FileFromString(
			corePassOverrideFilePath, "root", 0644, "")
		i.Config.Storage.Files = append(i.Config.Storage.Files, registryEnvFile)
	}

	// Add registry.env file
	registryEnvFile := ignasset.FileFromString(consts.RegistryEnvPath,
		"root", 0644, templates.GetRegistryEnv(consts.RegistryDataBootstrap))
	i.Config.Storage.Files = append(i.Config.Storage.Files, registryEnvFile)

	// Add public ssh key
	if applianceConfig.Config.SshKey != nil {
		passwdUser.SSHAuthorizedKeys = []igntypes.SSHAuthorizedKey{
			igntypes.SSHAuthorizedKey(*applianceConfig.Config.SshKey),
		}
	}
	i.Config.Passwd.Users = append(i.Config.Passwd.Users, passwdUser)

	err = addExtraManifests(&i.Config, extraManifests)
	if err != nil {
		return err
	}

	logrus.Debug("Successfully generated bootstrap ignition")

	return nil
}

// addExtraManifests is a non-exportable function copy-over from openshift/installer/pkg/asset/agent/image/ignition.go
func addExtraManifests(config *igntypes.Config, extraManifests *agentManifests.ExtraManifests) error {

	user := "root"
	mode := 0644

	config.Storage.Directories = append(config.Storage.Directories, igntypes.Directory{
		Node: igntypes.Node{
			Path: extraManifestPath,
			User: igntypes.NodeUser{
				Name: &user,
			},
			Overwrite: util.BoolToPtr(true),
		},
		DirectoryEmbedded1: igntypes.DirectoryEmbedded1{
			Mode: &mode,
		},
	})

	for _, file := range extraManifests.FileList {

		type unstructured map[string]interface{}

		yamlList, err := agentManifests.GetMultipleYamls[unstructured](file.Data)
		if err != nil {
			return errors.Wrapf(err, "could not decode YAML for %s", file.Filename)
		}

		for n, manifest := range yamlList {
			m, err := yaml.Marshal(manifest)
			if err != nil {
				return err
			}

			base := filepath.Base(file.Filename)
			ext := filepath.Ext(file.Filename)
			baseWithoutExt := strings.TrimSuffix(base, ext)
			baseFileName := filepath.Join(extraManifestPath, baseWithoutExt)
			fileName := fmt.Sprintf("%s-%d%s", baseFileName, n, ext)

			extraFile := ignition.FileFromBytes(fileName, user, mode, m)
			config.Storage.Files = append(config.Storage.Files, extraFile)
		}
	}

	return nil
}
