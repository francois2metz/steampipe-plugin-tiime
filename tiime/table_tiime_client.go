package tiime

import (
	"context"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeClient() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_client",
		Description: "A Client.",
		List: &plugin.ListConfig{
			Hydrate: listClient,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique id of the client."},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "Slug of the client."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the client."},
			{Name: "address", Type: proto.ColumnType_STRING, Description: "Address of the client."},
			{Name: "city", Type: proto.ColumnType_STRING, Description: "City of the client."},
			{Name: "payment_status", Type: proto.ColumnType_STRING, Description: "Payment status of the client."},
			{Name: "balance_including_taxes", Type: proto.ColumnType_DOUBLE, Description: "Balance of the client."},
			{Name: "billed_including_taxes", Type: proto.ColumnType_DOUBLE, Description: "Balance of the client."},
			{Name: "billed_excluding_taxes", Type: proto.ColumnType_DOUBLE, Description: "Balance of the client."},
		},
	}
}

func listClient(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_client.listClient", "connection_error", err)
		return nil, err
	}
	increment := 100
	opts := tiime.PaginationOpts{Start: 0, End: increment}

	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < int64(opts.End) {
		opts.End = int(*d.QueryContext.Limit)
	}
	for {
		clients, pagination, err := client.GetClients(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("tiime_client.listClient", err)
			return nil, err
		}
		for _, client := range clients {
			d.StreamListItem(ctx, client)
		}
		if pagination.Max != "*" {
			break
		}
		opts.Start = pagination.CurrentStart + increment + 1
		opts.End = pagination.CurrentEnd + increment + 1
		if d.RowsRemaining(ctx) <= 0 {
			break
		}
	}
	return nil, nil
}
