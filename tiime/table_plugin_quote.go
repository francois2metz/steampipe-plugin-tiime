package tiime

import (
	"context"
	"strings"
	"time"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeQuote() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_quote",
		Description: "A qute.",
		List: &plugin.ListConfig{
			Hydrate: listQuote,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "emission_date", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getQuote,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the quote.",
			},
			{
				Name:        "client_id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the client.",
			},
			{
				Name:        "compiled_number",
				Type:        proto.ColumnType_STRING,
				Description: "Unique number of the quote.",
			},
			{
				Name:        "number",
				Type:        proto.ColumnType_INT,
				Description: "Sequence number of the quote.",
			},
			{
				Name:        "emission_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Emission date of the qute.",
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_STRING,
				Description: "Template of the quote.",
			},
			{
				Name:        "color",
				Type:        proto.ColumnType_STRING,
				Description: "Color of the quote.",
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
				Description: "Quote comment.",
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the quote.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the quote (saved, accepted, refused, cancelled).",
			},
			{
				Name:        "lines",
				Type:        proto.ColumnType_JSON,
				Description: "Lines of the quote.",
				Hydrate:     getQuoteInfo,
			},
		},
	}
}

func listQuote(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_quote.listQuote", "connection_error", err)
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
		quotes, pagination, err := client.GetQuotes(ctx, opts, paginationOpts)
		if err != nil {
			plugin.Logger(ctx).Error("tiime_quote.listQuote", err)
			return nil, err
		}
		for _, quote := range quotes {
			d.StreamListItem(ctx, quote)
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

func getQuote(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()

	return getQuoteById(ctx, d, id)
}

func getQuoteInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quote := h.Item.(tiime.Quote)
	if quote.Lines != nil {
		return quote, nil
	}
	return getQuoteById(ctx, d, quote.ID)
}

func getQuoteById(ctx context.Context, d *plugin.QueryData, id int64) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_quote.getQuoteById", "connection_error", err)
		return nil, err
	}
	result, err := client.GetQuote(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_quote.getQuoteById", err)
		return nil, err
	}

	return result, nil
}
