package tiime

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func tableTiimeLabel() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_label",
		Description: "Label on transactions.",
		List: &plugin.ListConfig{
			Hydrate: listLabel,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the label.",
			},
			{
				Name:        "disabled",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "invoice_client",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "predicable",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "is_standard",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "label",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
		},
	}
}

func listLabel(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_label.listLabel", "connection_error", err)
		return nil, err
	}
	company_id, err := defaultCompanyID(d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_label.listLabel", "company error", err)
		return nil, err
	}
	labels, err := client.GetLabels(ctx, company_id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_label.listLabel", err)
		return nil, err
	}
	for _, label := range labels {
		d.StreamListItem(ctx, label)
	}
	return nil, nil
}
