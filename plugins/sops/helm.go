package sops

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HelmCLI() schema.Executable {
	return schema.Executable{
		Name:    "Helm with SOPS Secrets and Kubernetes",
		Runs:    []string{"helm"},
		DocsURL: sdk.URL("https://github.com/jkroepke/helm-secrets"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:   sdk.CredentialName("Kubeconfig"),
				Plugin: "kubernetes",
				// Kubeconfig - provisions KUBECONFIG env var pointing to temp file
			},
			{
				Name:     credname.SecretKey,
				Optional: true,
				// SOPS age key - provisions SOPS_AGE_KEY env var (optional, only for helm-secrets)
			},
		},
	}
}
