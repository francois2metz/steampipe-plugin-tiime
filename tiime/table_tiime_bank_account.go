package tiime

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeBankAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_bank_account",
		Description: "Bank accounts associated to the company",
		List: &plugin.ListConfig{
			Hydrate: listBankAccount,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique id of the bank account.",
			},
			{
				Name:        "synchronization_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "last_push_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "short_bank_name",
				Type:        proto.ColumnType_STRING,
				Description: ".",
			},
			{
				Name:        "bank_name",
				Type:        proto.ColumnType_STRING,
				Description: ".",
			},
			{
				Name:        "disableable",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "authorized_balance",
				Type:        proto.ColumnType_DOUBLE,
				Description: "",
			},
			{
				Name:        "pending_balance",
				Type:        proto.ColumnType_DOUBLE,
				Description: "",
			},
			{
				Name:        "closurable",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "iban",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "bic",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "balance_amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: ".",
			},
			{
				Name:        "balance_currency",
				Type:        proto.ColumnType_STRING,
				Description: ".",
			},
			{
				Name:        "balance_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "",
			},
			{
				Name:        "closed",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
			{
				Name:        "is_wallet",
				Type:        proto.ColumnType_BOOL,
				Description: "",
			},
		},
	}
}

func listBankAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_account.listBankAccount", "connection_error", err)
		return nil, err
	}
	company_id, err := defaultCompanyID(d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_account.listBankAccount", "company error", err)
		return nil, err
	}
	bankAccounts, err := client.GetBankAccounts(ctx, company_id)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_bank_account.listBankAccount", err)
		return nil, err
	}
	for _, bankAccount := range bankAccounts {
		d.StreamListItem(ctx, bankAccount)
	}

	return nil, nil
}
