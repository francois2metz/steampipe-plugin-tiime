package tiime

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/schema"
)

type tiimeConfig struct {
	Email     *string `cty:"email"`
	Password  *string `cty:"password"`
	CompanyID *int    `cty:"company_id"`
}

var ConfigSchema = map[string]*schema.Attribute{
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
