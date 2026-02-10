package helm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "helm",
		Platform: schema.PlatformInfo{
			Name:     "Helm",
			Homepage: sdk.URL("https://helm.sh"),
		},
		Credentials: []schema.CredentialType{
			HelmCredentials(),
		},
		Executables: []schema.Executable{
			HelmCLI(),
			HelmfileCLI(),
		},
	}
}
