package helm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HelmCLI() schema.Executable {
	return schema.Executable{
		Name:    "Helm",
		Runs:    []string{"helm"},
		DocsURL: sdk.URL("https://helm.sh/docs/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{Name: sdk.CredentialName("Kubeconfig")},
			{Name: credname.SecretKey, Plugin: "sops", Optional: true},
			{Name: credname.PersonalAccessToken, Plugin: "digitalocean", Optional: true},
		},
	}
}
