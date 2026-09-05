package tiime

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/schema"
)

type tiimeConfig struct {
	ClientID  *string `cty:"client_id"`
	Email     *string `cty:"email"`
	Password  *string `cty:"password"`
	CompanyID *int    `cty:"company_id"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"client_id": {
		Type: schema.TypeString,
	},
	"email": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"company_id": {
		Type: schema.TypeInt,
	},
}

func ConfigInstance() interface{} {
	return &tiimeConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) tiimeConfig {
	if connection == nil || connection.GetConfig() == nil {
		return tiimeConfig{}
	}
	config, _ := connection.GetConfig().(tiimeConfig)
	return config
}
