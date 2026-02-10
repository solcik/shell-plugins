package helm

import (
	"context"
	"encoding/base64"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func HelmCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    sdk.CredentialName("Helm Credentials"),
		DocsURL: sdk.URL("https://helm.sh/docs/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Credential,
				MarkdownDescription: "Base64-encoded kubeconfig YAML file contents.",
				Secret:              true,
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Age secret key used by SOPS for encryption and decryption.",
				Secret:              true,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Prefix: "AGE-SECRET-KEY-",
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: &helmCredentialsProvisioner{},
		Importer: importer.TryAll(
			TryKubeconfigFile(),
			TryAgeKeyEnvVar(),
		),
	}
}

func TryKubeconfigFile() sdk.Importer {
	return importer.TryFile("~/.kube/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		encoded := base64.StdEncoding.EncodeToString(contents)
		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Credential: encoded,
			},
			NameHint: importer.SanitizeNameHint("default"),
		})
	})
}

func TryAgeKeyEnvVar() sdk.Importer {
	return importer.TryEnvVarPair(map[string]sdk.FieldName{
		"SOPS_AGE_KEY": fieldname.PrivateKey,
	})
}
