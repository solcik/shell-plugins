package helm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
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
			{
				Name: sdk.CredentialName("Helm Credentials"),
			},
		},
	}
}
