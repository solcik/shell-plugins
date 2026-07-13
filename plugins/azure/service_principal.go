package azure

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ServicePrincipal() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.ServicePrincipal,
		DocsURL:       sdk.URL("https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli-service-principal"),
		ManagementURL: sdk.URL("https://entra.microsoft.com/#view/Microsoft_AAD_RegisteredApps/ApplicationsListBlade"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ClientID,
				MarkdownDescription: "Application (client) ID of the service principal.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						// Client IDs are lowercased UUIDs.
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.ClientSecret,
				MarkdownDescription: "Secret created for the service principal.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.TenantID,
				MarkdownDescription: "Tenant ID which will be authenticated to.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						// Tenant IDs are lowercased UUIDs.
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.Subscription,
				MarkdownDescription: "Subscription ID or name to select as the default subscription.",
				Secret:              false,
				Optional:            true,
			},
		},
		DefaultProvisioner: &azureConfigDirProvisioner{},
	}
}
