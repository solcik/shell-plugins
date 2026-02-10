package helm

import (
	"encoding/base64"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestKubeconfigProvisioner(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	plugintest.TestProvisioner(t, Kubeconfig().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"kubeconfig": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: encodedConfig,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"KUBECONFIG": "/tmp/config",
				},
			},
		},
	})
}

func TestKubeconfigImporter(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	plugintest.TestImporter(t, Kubeconfig().Importer, map[string]plugintest.ImportCase{
		"kubeconfig file": {
			Files: map[string]string{
				"~/.kube/config": rawConfig,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Credential: encodedConfig,
					},
					NameHint: "",
				},
			},
		},
	})
}
