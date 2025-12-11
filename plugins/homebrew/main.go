// Package main implements the Homebrew formula publishing plugin for ReleasePilot.
package main

import (
	"github.com/felixgeelhaar/release-pilot/pkg/plugin"
)

func main() {
	plugin.Serve(&HomebrewPlugin{})
}
