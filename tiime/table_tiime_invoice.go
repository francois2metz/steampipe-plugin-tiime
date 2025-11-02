package tiime

import (
	"context"
	"strings"
	"time"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeInvoice() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_invoice",
		Description: "An invoice.",
		List: &plugin.ListConfig{
			Hydrate: listInvoice,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "emission_date", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getInvoice,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the invoice.",
			},
			{
				Name:        "client_id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the client.",
			},
			{
				Name:        "compiled_number",
				Type:        proto.ColumnType_STRING,
				Description: "Unique number of the invoice.",
			},
			{
				Name:        "number",
				Type:        proto.ColumnType_INT,
				Description: "Sequence number of the invoice.",
			},
			{
				Name:        "emission_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Emission date of the invoice.",
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_STRING,
				Description: "Template of the invoice.",
			},
			{
				Name:        "color",
				Type:        proto.ColumnType_STRING,
				Description: "Color of the invoice.",
			},
			{
				Name:        "client_name",
				Type:        proto.ColumnType_STRING,
				Description: "Client name.",
			},
			{
				Name:        "total_excluding_taxes",
				Type:        proto.ColumnType_DOUBLE,
				Description: "Total amount excluding taxes.",
			},
			{
				Name:        "total_including_taxes",
				Type:        proto.ColumnType_DOUBLE,
				Description: "Total amount including taxes.",
			},
			{
				Name:        "comment",
				Type:        proto.ColumnType_STRING,
				Description: "Invoice comment.",
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the invoice.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the invoice (draft, sent, paid, ...).",
			},
			{
				Name:        "lines",
				Type:        proto.ColumnType_JSON,
				Description: "Lines of the invoice.",
				Hydrate:     getInvoiceInfo,
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags of the invoice.",
				Hydrate:     getInvoiceInfo,
			},
		},
	}
}

func listInvoice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.listInvoice", "connection_error", err)
		return nil, err
	}
	maxItem := 100
	paginationOpts := tiime.PaginationOpts{Start: 0, End: maxItem - 1}

	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < int64(paginationOpts.End) {
		paginationOpts.End = int(*d.QueryContext.Limit)
	}
	status := d.EqualsQuals["status"].GetStringValue()
	emission_date := d.Quals["emission_date"]
	opts := tiime.ListQueryOpts{
		Status: status,
	}
	if emission_date != nil {
		var date_query []string
		for _, q := range emission_date.Quals {
			if q.Value.GetTimestampValue() != nil {
				date := q.Value.GetTimestampValue().AsTime().Format(time.DateOnly)
				if q.Operator == "=" {
					date_query = append(date_query, date)
				} else {
					date_query = append(date_query, q.Operator+date)
				}
			}
		}
		opts.EmissionDate = strings.Join(date_query, ",")
	}
	for {
		invoices, pagination, err := client.GetInvoices(ctx, opts, paginationOpts)
		if err != nil {
			plugin.Logger(ctx).Error("tiime_invoice.listInvoice", err)
			return nil, err
		}
		for _, invoice := range invoices {
			d.StreamListItem(ctx, invoice)
		}
		if pagination.Max != "*" {
			break
		}
		paginationOpts.Start += maxItem
		paginationOpts.End += maxItem
		if d.RowsRemaining(ctx) <= 0 {
			break
		}
	}
	return nil, nil
}

func getInvoice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()

	return getInvoiceById(ctx, d, id)
}

func getInvoiceInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	invoice := h.Item.(tiime.Invoice)
	if invoice.Lines != nil {
		return invoice, nil
	}
	return getInvoiceById(ctx, d, invoice.ID)
}

func getInvoiceById(ctx context.Context, d *plugin.QueryData, id int64) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.getInvoiceById", "connection_error", err)
		return nil, err
	}
	result, err := client.GetInvoice(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_invoice.getInvoiceById", err)
		return nil, err
	}

	return result, nil
}
