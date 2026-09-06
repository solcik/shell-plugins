package wrangler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

// AlchemyCLI is Alchemy, an infrastructure-as-code tool that deploys
// Cloudflare Workers. It authenticates with the same account ID and API
// token as Wrangler, so it reuses this plugin's credential.
//
// Alchemy reads CLOUDFLARE_ACCOUNT_ID and CLOUDFLARE_API_TOKEN only when
// its auth profile records `method: "env"`. Configure the profile with
// `alchemy login --configure` and select "Environment Variables".
func AlchemyCLI() schema.Executable {
	return schema.Executable{
		Name:      "Alchemy",
		Runs:      []string{"alchemy"},
		DocsURL:   sdk.URL("https://alchemy.run/getting-started/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
