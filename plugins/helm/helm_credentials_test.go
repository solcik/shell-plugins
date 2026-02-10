package helm

import (
	"encoding/base64"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestHelmCredentialsProvisioner(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	plugintest.TestProvisioner(t, HelmCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"kubeconfig only": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: encodedConfig,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"KUBECONFIG": "/tmp/config",
				},
			},
		},
		"kubeconfig and sops age key": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: encodedConfig,
				fieldname.PrivateKey: "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"KUBECONFIG":   "/tmp/config",
					"SOPS_AGE_KEY": "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
				},
			},
		},
	})
}

func TestHelmCredentialsImporter(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	plugintest.TestImporter(t, HelmCredentials().Importer, map[string]plugintest.ImportCase{
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
		"age key environment": {
			Environment: map[string]string{
				"SOPS_AGE_KEY": "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.PrivateKey: "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
					},
				},
			},
		},
	})
}
