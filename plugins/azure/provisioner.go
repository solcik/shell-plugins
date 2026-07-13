package azure

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// azureConfigDirProvisioner authenticates the Azure CLI by running
// `az login --service-principal` against a temporary AZURE_CONFIG_DIR.
// The az CLI has no environment-variable credential support, so login state
// must live in a config dir; using a throwaway one keeps the user's ~/.azure
// untouched and leaves no credentials behind outside the plugin temp dir.
type azureConfigDirProvisioner struct{}

func (p *azureConfigDirProvisioner) Description() string {
	return "Provision a temporary AZURE_CONFIG_DIR logged in with a service principal"
}

func (p *azureConfigDirProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	configDir := filepath.Join(in.TempDir, "azure")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		out.AddError(err)
		return
	}
	env := append(os.Environ(), "AZURE_CONFIG_DIR="+configDir)

	login := exec.CommandContext(ctx, "az", "login", "--service-principal",
		"--username", in.ItemFields[fieldname.ClientID],
		"--password", in.ItemFields[fieldname.ClientSecret],
		"--tenant", in.ItemFields[fieldname.TenantID],
		"--output", "none", "--only-show-errors")
	login.Env = env
	if output, err := login.CombinedOutput(); err != nil {
		out.AddError(fmt.Errorf("az login failed: %w: %s", err, output))
		return
	}

	if subscription := in.ItemFields[fieldname.Subscription]; subscription != "" {
		set := exec.CommandContext(ctx, "az", "account", "set",
			"--subscription", subscription, "--only-show-errors")
		set.Env = env
		if output, err := set.CombinedOutput(); err != nil {
			out.AddError(fmt.Errorf("az account set failed: %w: %s", err, output))
			return
		}
	}

	out.AddEnvVar("AZURE_CONFIG_DIR", configDir)
}

func (p *azureConfigDirProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Best effort: op doesn't call Deprovision for local plugins; the temp dir
	// is cleaned up externally (tmpfiles rule on /tmp/1PasswordShellPlugins-*).
	os.RemoveAll(filepath.Join(in.TempDir, "azure"))
}
