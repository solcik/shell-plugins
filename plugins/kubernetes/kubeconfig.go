package kubernetes

import (
	"context"
	"encoding/base64"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Kubeconfig() schema.CredentialType {
	return schema.CredentialType{
		Name:    sdk.CredentialName("Kubeconfig"),
		DocsURL: sdk.URL("https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Credential,
				MarkdownDescription: "Base64-encoded kubeconfig YAML file contents.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.TempFile(
			kubeconfigFileContents,
			provision.Filename("config"),
			provision.SetPathAsEnvVar("KUBECONFIG"),
		),
		Importer: importer.TryAll(
			TryKubeconfigFile(),
		),
	}
}

func kubeconfigFileContents(in sdk.ProvisionInput) ([]byte, error) {
	encoded := in.ItemFields[fieldname.Credential]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func TryKubeconfigFile() sdk.Importer {
	return importer.TryFile("~/.kube/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// Import existing kubeconfig as base64-encoded
		encoded := base64.StdEncoding.EncodeToString(contents)
		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Credential: encoded,
			},
			NameHint: importer.SanitizeNameHint("default"),
		})
	})
}
