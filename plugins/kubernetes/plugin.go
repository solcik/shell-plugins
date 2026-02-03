package kubernetes

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "kubernetes",
		Platform: schema.PlatformInfo{
			Name:     "Kubernetes",
			Homepage: sdk.URL("https://kubernetes.io"),
		},
		Credentials: []schema.CredentialType{
			Kubeconfig(),
		},
		Executables: []schema.Executable{
			KubectlCLI(),
			SternCLI(),
			K9sCLI(),
		},
	}
}
