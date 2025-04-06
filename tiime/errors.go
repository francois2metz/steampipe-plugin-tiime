package tiime

import (
	"context"
	"net/http"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func shouldIgnoreError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
	requestError, ok := err.(*tiime.APIError)
	if !ok {
		return false
	}

	if requestError.ErrorStatus == http.StatusNotFound {
		return true
	}
	return false
}
