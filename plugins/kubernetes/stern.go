package kubernetes

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func SternCLI() schema.Executable {
	return schema.Executable{
		Name:    "stern",
		Runs:    []string{"stern"},
		DocsURL: sdk.URL("https://github.com/stern/stern"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: sdk.CredentialName("Kubeconfig"),
			},
		},
	}
}
