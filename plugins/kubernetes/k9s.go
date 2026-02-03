package kubernetes

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func K9sCLI() schema.Executable {
	return schema.Executable{
		Name:    "k9s",
		Runs:    []string{"k9s"},
		DocsURL: sdk.URL("https://k9scli.io/"),
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
