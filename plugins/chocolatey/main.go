// Package main implements the Chocolatey plugin for ReleasePilot.
package main

import (
	"github.com/felixgeelhaar/release-pilot/pkg/plugin"
)

func main() {
	plugin.Serve(&ChocolateyPlugin{})
}
