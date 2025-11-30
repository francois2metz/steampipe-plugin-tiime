package tiime

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTiimeCompany() *plugin.Table {
	return &plugin.Table{
		Name:        "tiime_company",
		Description: "The companies handled by Tiime.",
		List: &plugin.ListConfig{
			Hydrate: listCompany,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique id of the company."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the company."},
			{Name: "siret", Type: proto.ColumnType_STRING, Description: "Siret of the company."},
			{Name: "street", Type: proto.ColumnType_STRING, Description: "Street address of the company."},
			{Name: "city", Type: proto.ColumnType_STRING, Description: "City of the company."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country of the company."},
			{Name: "legal_form", Type: proto.ColumnType_STRING, Description: "Legal form of the company."},
		},
	}
}

func listCompany(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_company.listCompany", "connection_error", err)
		return nil, err
	}
	companies, err := client.GetCompanies(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("tiime_company.listCompany", err)
		return nil, err
	}
	for _, company := range companies {
		d.StreamListItem(ctx, company)
	}
	return nil, nil
}
