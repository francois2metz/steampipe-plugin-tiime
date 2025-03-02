package tiime

import (
	"context"
	"errors"
	"os"
	"strconv"

	tiime "github.com/francois2metz/steampipe-plugin-tiime/tiime/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*tiime.Client, error) {
	// get tiime client from cache
	cacheKey := "tiime"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*tiime.Client), nil
	}

	tiimeConfig := GetConfig(d.Connection)

	email := os.Getenv("TIIME_EMAIL")
	password := os.Getenv("TIIME_PASSWORD")
	var company_id int
	if os.Getenv("TIIME_COMPANY_ID") != "" {
		var err error
		company_id, err = strconv.Atoi(os.Getenv("TIIME_COMPANY_ID"))
		if err != nil {
			return nil, errors.New("TIIME_COMPANY_ID environnement variable error")
		}
	}

	if tiimeConfig.Email != nil {
		email = *tiimeConfig.Email
	}
	if tiimeConfig.Password != nil {
		password = *tiimeConfig.Password
	}
	if tiimeConfig.CompanyID != nil {
		company_id = *tiimeConfig.CompanyID
	}

	if email == "" {
		return nil, errors.New("'email' must be set in the connection configuration. Edit your connection configuration file or set the TIIME_EMAIL environment variable and then restart Steampipe")
	}

	if password == "" {
		return nil, errors.New("'password' must be set in the connection configuration. Edit your connection configuration file or set the TIIME_PASSWORD environment variable and then restart Steampipe")
	}

	if company_id == 0 {
		return nil, errors.New("'company_id' must be set in the connection configuration. Edit your connection configuration file or set the TIIME_COMPANY_ID environment variable and then restart Steampipe")
	}

	config := tiime.ClientConfig{
		Email:     email,
		Password:  password,
		CompanyID: company_id,
	}
	client, err := tiime.New(ctx, config)
	if err != nil {
		return nil, err
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
