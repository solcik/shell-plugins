package kubernetes

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func KubectlCLI() schema.Executable {
	return schema.Executable{
		Name:    "kubectl",
		Runs:    []string{"kubectl"},
		DocsURL: sdk.URL("https://kubernetes.io/docs/reference/kubectl/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: sdk.CredentialName("Kubeconfig"),
			},
		},
	}
}
