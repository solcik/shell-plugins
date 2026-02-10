package helm

import (
	"context"
	"encoding/base64"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type helmCredentialsProvisioner struct{}

func (p *helmCredentialsProvisioner) Description() string {
	return "Provision kubeconfig file and optional SOPS age key for Helm"
}

func (p *helmCredentialsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// Decode base64 kubeconfig and write to temp file
	encoded := in.ItemFields[fieldname.Credential]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		out.AddError(err)
		return
	}

	configPath := filepath.Join(in.TempDir, "config")
	out.AddSecretFile(configPath, decoded)
	out.AddEnvVar("KUBECONFIG", configPath)

	// Optionally provision SOPS age key
	if ageKey, ok := in.ItemFields[fieldname.PrivateKey]; ok && ageKey != "" {
		out.AddEnvVar("SOPS_AGE_KEY", ageKey)
	}
}

func (p *helmCredentialsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Temp files are automatically cleaned up
}
