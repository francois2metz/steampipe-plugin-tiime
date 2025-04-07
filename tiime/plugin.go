package tiime

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-tiime",
		DefaultTransform: transform.FromGo(),
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreError,
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"tiime_bank_account": tableTiimeBankAccount(),
			"tiime_client":       tableTiimeClient(),
			"tiime_invoice":      tableTiimeInvoice(),
			"tiime_quote":        tableTiimeQuote(),
		},
	}
	return p
}
