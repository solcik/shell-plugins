package helm

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type helmKubeconfigProvisioner struct{}

func (p *helmKubeconfigProvisioner) Description() string {
	return "Provision kubeconfig file for Helm"
}

func (p *helmKubeconfigProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// Decode base64 kubeconfig
	encoded := in.ItemFields[fieldname.Credential]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		out.AddError(err)
		return
	}

	// Write kubeconfig as a real file (not via out.AddSecretFile which creates a FIFO).
	// Helm reads the kubeconfig multiple times, and FIFOs block on the second read.
	configPath := filepath.Join(in.TempDir, "config")
	if err := os.WriteFile(configPath, decoded, 0600); err != nil {
		out.AddError(err)
		return
	}
	out.AddEnvVar("KUBECONFIG", configPath)
}

func (p *helmKubeconfigProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Remove kubeconfig written directly to disk
	configPath := filepath.Join(in.TempDir, "config")
	os.Remove(configPath)
}
