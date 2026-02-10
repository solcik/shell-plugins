package helm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HelmfileCLI() schema.Executable {
	return schema.Executable{
		Name:    "Helmfile",
		Runs:    []string{"helmfile"},
		DocsURL: sdk.URL("https://github.com/helmfile/helmfile"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{Name: sdk.CredentialName("Kubeconfig")},
			{Name: credname.SecretKey, Plugin: "sops", Optional: true},
		},
	}
}
