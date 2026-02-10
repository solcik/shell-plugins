package helm

import (
	"context"
	"encoding/base64"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Kubeconfig() schema.CredentialType {
	return schema.CredentialType{
		Name:    sdk.CredentialName("Kubeconfig"),
		DocsURL: sdk.URL("https://helm.sh/docs/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Credential,
				MarkdownDescription: "Base64-encoded kubeconfig YAML file contents.",
				Secret:              true,
			},
		},
		DefaultProvisioner: &helmKubeconfigProvisioner{},
		Importer: importer.TryAll(
			TryKubeconfigFile(),
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
