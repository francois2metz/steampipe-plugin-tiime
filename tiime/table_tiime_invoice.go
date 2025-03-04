package tiime

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeInvoice() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_invoice",
		Description: "An invoice.",
		List: &plugin.ListConfig{
			Hydrate: listInvoice,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getInvoice,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique id of the invoice."},
			{Name: "client_id", Type: proto.ColumnType_INT, Description: "Unique id of the client."},
			{Name: "compiled_number", Type: proto.ColumnType_STRING, Description: "Unique number of the invoice."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "Sequence number of the invoice."},
			{Name: "emission_date", Type: proto.ColumnType_STRING, Description: "Emission date of the invoice."},
			{Name: "template", Type: proto.ColumnType_STRING, Description: "Template of the invoice."},
			{Name: "color", Type: proto.ColumnType_STRING, Description: "Color of the invoice."},
			{Name: "client_name", Type: proto.ColumnType_STRING, Description: "Client name."},
			{Name: "total_excluding_taxes", Type: proto.ColumnType_DOUBLE, Description: "Total amount excluding taxes."},
			{Name: "total_including_taxes", Type: proto.ColumnType_DOUBLE, Description: "Total amount including taxes."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Comment."},
		},
	}
}

func listInvoice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.listInvoice", "connection_error", err)
		return nil, err
	}
	invoices, err := client.GetInvoices(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.listInvoice", err)
		return nil, err
	}
	for _, invoice := range invoices {
		d.StreamListItem(ctx, invoice)
	}
	return nil, nil
}

func getInvoice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.getInvoice", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["id"].GetInt64Value()
	result, err := client.GetInvoice(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.getInvoice", err)
		return nil, err
	}

	return result, nil
}
