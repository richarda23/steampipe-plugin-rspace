package main

import (
    "github.com/turbot/steampipe-plugin-sdk/plugin"
    "github.com/richarda23/steampipe-plugin-rspace/rspace"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{PluginFunc: rspace.Plugin})
}
