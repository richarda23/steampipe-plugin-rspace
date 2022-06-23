package main

import (
	"github.com/richarda23/steampipe-plugin-rspace/rspace"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: rspace.Plugin})
}
