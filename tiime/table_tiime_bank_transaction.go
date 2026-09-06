package tiime

import (
	"context"
	"strings"
	"time"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func tableTiimeBankTransaction() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_bank_transaction",
		Description: "A bank transaction.",
		List: &plugin.ListConfig{
			Hydrate: listBankTransaction,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "transaction_date", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBankTransaction,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the bank transaction.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the transaction (done, failed).",
			},
			{
				Name:        "transaction_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "realization_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "vat_application_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: "",
			},
			{
				Name:        "currency",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "comment",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
		},
	}
}

func listBankTransaction(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_transaction.listBankTransaction", "connection_error", err)
		return nil, err
	}
	company_id, err := defaultCompanyID(d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_transaction.listBankTransaction", "company error", err)
		return nil, err
	}
	maxItem := 100
	paginationOpts := tiime.PaginationOpts{Start: 0, End: maxItem - 1}

	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < int64(paginationOpts.End) {
		paginationOpts.End = int(*d.QueryContext.Limit)
	}
	transaction_date := d.Quals["transaction_date"]
	opts := tiime.ListTransactionOpts{}
	if transaction_date != nil {
		var date_query []string
		for _, q := range transaction_date.Quals {
			if q.Value.GetTimestampValue() != nil {
				date := q.Value.GetTimestampValue().AsTime().Format(time.DateOnly)
				if q.Operator == "=" {
					date_query = append(date_query, date)
				} else {
					date_query = append(date_query, q.Operator+date)
				}
			}
		}
		opts.TransactionDate = strings.Join(date_query, ",")
	}
	for {
		transactions, pagination, err := client.GetBankTransactions(ctx, company_id, opts, paginationOpts)
		if err != nil {
			plugin.Logger(ctx).Error("tiime_bank_transaction.listBankTransaction", err)
			return nil, err
		}
		for _, transaction := range transactions.Transactions {
			d.StreamListItem(ctx, transaction)
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

func getBankTransaction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_transaction.getBankTransaction", "connection_error", err)
		return nil, err
	}
	company_id, err := defaultCompanyID(d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_transaction.getBankTransaction", "company error", err)
		return nil, err
	}
	id := d.EqualsQuals["id"].GetInt64Value()
	result, err := client.GetBankTransaction(ctx, company_id, id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_transaction.getBankTransaction", err)
		return nil, err
	}

	return result, nil
}
