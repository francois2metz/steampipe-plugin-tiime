package main

import (
	"github.com/francois2metz/steampipe-plugin-tiime/tiime"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: tiime.Plugin})
}
